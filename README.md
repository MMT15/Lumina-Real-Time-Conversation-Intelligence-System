# Lumina: Real-Time Conversation Intelligence System

Lumina is a "Cresta-level" conversation intelligence system designed for high-performance backend processing, semantic search, and automated task suggestions.

## Key Features

- **Auto-Task Suggestion**: Automatically generates tickets based on conversation context (e.g., "I want a refund" -> "Pending" ticket).
- **Semantic Search**: Powered by Elasticsearch to find conversations where customers were unhappy with pricing or other specific themes.
- **Sentiment Analysis**: Detects negative sentiment and instantly alerts managers.

## Tech Stack

- **Backend**: Go (Golang) with Gin for REST and gRPC for high performance.
- **Database**: PostgreSQL (GORM) for relational data.
- **Search Engine**: Elasticsearch for complex semantic queries.
- **Containerization**: Docker & Docker Compose.
- **AI/LLM**: Ollama (Llama 3.2:3b) for local sentiment analysis and intent detection.

## User Story

"Utilizatorul 'Ion' trimite un mesaj: 'Sunt foarte supărat că pachetul nu a ajuns!'. Sistemul Lumina detectează instant sentimentul negativ, trimite o alertă managerului și creează un task de urmărire a coletului în baza de date, totul în sub o secundă."

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/MMT15/Lumina-Real-Time-Conversation-Intelligence-System.git
   cd Lumina
   ```
2. Start the infrastructure:
   ```bash
   docker-compose up -d
   ```
3. Run the application:
   ```bash
   go run main.go
   ```
