package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/MMT15/Lumina/database"
	"github.com/MMT15/Lumina/models"
	"github.com/MMT15/Lumina/pb"
	"github.com/MMT15/Lumina/services"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedConversationServiceServer
}

func (s *server) ProcessMessage(ctx context.Context, in *pb.MessageRequest) (*pb.MessageResponse, error) {
	// AI Analysis with Ollama
	analysis, err := services.AnalyzeSentimentAndIntent(in.Content)
	if err != nil {
		log.Printf("Ollama analysis failed: %v", err)
		analysis = &services.AIAnalysis{Sentiment: "Neutral", CreateTicket: false}
	}

	// Create Conversation record
	conv := models.Conversation{
		UserID:    in.UserId,
		Message:   in.Content,
		Sentiment: analysis.Sentiment,
	}

	if err := database.DB.Create(&conv).Error; err != nil {
		return nil, err
	}

	// Auto-Ticket generation
	ticketCreated := false
	ticketStatus := ""
	if analysis.CreateTicket {
		ticket := models.Ticket{
			ConversationID: conv.ID,
			Description:    analysis.TicketSummary,
			Status:         models.StatusPending,
		}
		if err := database.DB.Create(&ticket).Error; err == nil {
			ticketCreated = true
			ticketStatus = string(models.StatusPending)
		}
	}

	// Index in Elasticsearch
	services.IndexConversation(conv)

	return &pb.MessageResponse{
		Sentiment:     analysis.Sentiment,
		TicketCreated: ticketCreated,
		TicketStatus:  ticketStatus,
	}, nil
}

func main() {
	// Initialize Database
	database.InitDB()

	// Initialize Elasticsearch
	services.InitElasticsearch()

	// Start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterConversationServiceServer(s, &server{})
		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Lumina Intelligence System (Go) is up and running!",
		})
	})

	r.POST("/analyze", func(c *gin.Context) {
		var input struct {
			UserID string `json:"user_id" binding:"required"`
			Text   string `json:"text" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// AI Analysis with Ollama
		analysis, err := services.AnalyzeSentimentAndIntent(input.Text)
		if err != nil {
			log.Printf("Ollama analysis failed: %v", err)
			// Fallback if AI fails
			analysis = &services.AIAnalysis{Sentiment: "Neutral", CreateTicket: false}
		}

		// Create Conversation record
		conv := models.Conversation{
			UserID:    input.UserID,
			Message:   input.Text,
			Sentiment: analysis.Sentiment,
		}

		if err := database.DB.Create(&conv).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save conversation"})
			return
		}

		// Auto-Ticket generation based on AI intent
		var ticketID uint
		if analysis.CreateTicket {
			ticket := models.Ticket{
				ConversationID: conv.ID,
				Description:    analysis.TicketSummary,
				Status:         models.StatusPending,
			}
			if err := database.DB.Create(&ticket).Error; err == nil {
				ticketID = ticket.ID
			}
		}

		// Index in Elasticsearch
		services.IndexConversation(conv)

		c.JSON(http.StatusOK, gin.H{
			"id":             conv.ID,
			"user_id":        conv.UserID,
			"message":        conv.Message,
			"sentiment":      conv.Sentiment,
			"intent":         analysis.Intent,
			"ticket_created": analysis.CreateTicket,
			"ticket_id":      ticketID,
			"status":         "Processed",
		})
	})

	r.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Search query 'q' is required"})
			return
		}

		results, err := services.SearchConversations(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"query":   query,
			"results": results,
		})
	})

	r.Run(":8080")
}
