package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type OllamaResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type AIAnalysis struct {
	Sentiment     string `json:"sentiment"`
	Score         float64 `json:"score"`
	Intent        string  `json:"intent"`
	CreateTicket  bool    `json:"create_ticket"`
	TicketSummary string  `json:"ticket_summary"`
}

func AnalyzeSentimentAndIntent(text string) (*AIAnalysis, error) {
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434/api/chat"
	}

	prompt := fmt.Sprintf(`Analyze the following customer message: "%s".
Return ONLY a JSON object with these keys: 
"sentiment" (string: Positive, Negative, Neutral), 
"score" (float 0-1), 
"intent" (string: Refund, Complaint, Inquiry, Praise), 
"create_ticket" (boolean: true if user is frustrated or wants action like refund/complaint), 
"ticket_summary" (string: short description for a ticket).`, text)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": "llama3.2:3b",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"stream": false,
	})

	resp, err := http.Post(ollamaURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, err
	}

	content := ollamaResp.Message.Content
	// Find the start and end of the JSON object
	start := -1
	end := -1
	for i, char := range content {
		if char == '{' {
			start = i
			break
		}
	}
	for i := len(content) - 1; i >= 0; i-- {
		if content[i] == '}' {
			end = i
			break
		}
	}

	if start == -1 || end == -1 {
		return nil, fmt.Errorf("no JSON object found in AI response: %s", content)
	}

	jsonStr := content[start : end+1]
	var analysis AIAnalysis
	err = json.Unmarshal([]byte(jsonStr), &analysis)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON: %v, content: %s", err, jsonStr)
	}

	return &analysis, nil
}
