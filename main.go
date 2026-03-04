package main

import (
	"net/http"

	"github.com/MMT15/Lumina/database"
	"github.com/MMT15/Lumina/models"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	database.InitDB()

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

		// Create Conversation record
		conv := models.Conversation{
			UserID:  input.UserID,
			Message: input.Text,
		}

		// TODO: Implement Ollama sentiment analysis
		// For now, mock it
		conv.Sentiment = "Neutral"

		if err := database.DB.Create(&conv).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save conversation"})
			return
		}

		// TODO: Implement Auto-Ticket generation based on intent
		// TODO: Implement Elasticsearch indexing

		c.JSON(http.StatusOK, gin.H{
			"id":        conv.ID,
			"user_id":   conv.UserID,
			"message":   conv.Message,
			"sentiment": conv.Sentiment,
			"status":    "Analysis in progress",
		})
	})

	r.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		// TODO: Implement Elasticsearch semantic search
		c.JSON(http.StatusOK, gin.H{
			"query":   query,
			"results": []string{},
		})
	})

	r.Run(":8080")
}
