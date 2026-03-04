package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/MMT15/Lumina/models"
	"github.com/olivere/elastic/v7"
)

var ESClient *elastic.Client

func InitElasticsearch() {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://localhost:9200"
	}

	client, err := elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false), // Important for local/docker
	)
	if err != nil {
		log.Printf("Elasticsearch connection failed: %v", err)
		return
	}

	ESClient = client
	fmt.Println("Elasticsearch connection established")

	// Ensure index exists
	exists, err := client.IndexExists("conversations").Do(context.Background())
	if err == nil && !exists {
		_, err = client.CreateIndex("conversations").Do(context.Background())
		if err != nil {
			log.Printf("Failed to create ES index: %v", err)
		}
	}
}

func IndexConversation(conv models.Conversation) {
	if ESClient == nil {
		return
	}

	_, err := ESClient.Index().
		Index("conversations").
		Id(fmt.Sprintf("%d", conv.ID)).
		BodyJson(conv).
		Do(context.Background())

	if err != nil {
		log.Printf("Failed to index conversation in ES: %v", err)
	}
}

func SearchConversations(queryText string) ([]models.Conversation, error) {
	if ESClient == nil {
		return nil, fmt.Errorf("ES client not initialized")
	}

	// Simple match query on message field
	// For "Cresta-level" semantic search, we would use embeddings or more complex queries
	// but here we use match query as a starting point.
	query := elastic.NewMatchQuery("message", queryText)

	searchResult, err := ESClient.Search().
		Index("conversations").
		Query(query).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	var results []models.Conversation
	for _, hit := range searchResult.Hits.Hits {
		var conv models.Conversation
		err := json.Unmarshal(hit.Source, &conv)
		if err != nil {
			continue
		}
		results = append(results, conv)
	}

	return results, nil
}
