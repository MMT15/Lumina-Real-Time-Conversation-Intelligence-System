package main

import (
	"context"
	"log"
	"time"

	"github.com/MMT15/Lumina/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewConversationServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r, err := c.ProcessMessage(ctx, &pb.MessageRequest{
		UserId:  "maria456",
		Content: "Mâncarea a fost groaznică și rece!",
	})
	if err != nil {
		log.Fatalf("could not process message: %v", err)
	}

	log.Printf("Sentiment: %s", r.GetSentiment())
	log.Printf("Ticket Created: %v", r.GetTicketCreated())
	log.Printf("Ticket Status: %s", r.GetTicketStatus())
}
