package main

import (
	"log"
	"net/http"

	"github.com/MMT15/Lumina/database"
	"github.com/MMT15/Lumina/models"
	"github.com/MMT15/Lumina/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	database.InitDB()

	// Initialize Elasticsearch
	services.InitElasticsearch()

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
