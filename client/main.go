package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc-learn/proto"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewPaymentsServiceClient(conn)
	req := &pb.BalanceRequest{UserId: 1}
	resp, err := client.GetBalance(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling GetBalance: %v", err)
	}
	log.Printf("Result: %v", resp)
}
