package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Testing Lumina System...")

	// 1. Test /analyze with a complaint
	fmt.Println("\n1. Testing /analyze (Complaint/Refund)...")
	message := map[string]string{
		"user_id": "ion123",
		"text":    "Sunt foarte supărat că pachetul nu a ajuns! Vreau o rambursare imediată!",
	}
	body, _ := json.Marshal(message)
	resp, err := http.Post("http://localhost:8080/analyze", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error: %v (Make sure the server is running on :8080)\n", err)
	} else {
		defer resp.Body.Close()
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		fmt.Printf("Response: %+v\n", result)
	}

	time.Sleep(2 * time.Second) // Wait for ES indexing

	// 2. Test /search
	fmt.Println("\n2. Testing /search (Negative Sentiment/Complaint)...")
	resp, err = http.Get("http://localhost:8080/search?q=supărat")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		defer resp.Body.Close()
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		fmt.Printf("Search Results: %+v\n", result)
	}
}
