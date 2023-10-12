package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc-learn/proto"
	"io"
	"log"
)

func getBalance(client pb.PaymentsServiceClient) {
	req := &pb.BalanceRequest{UserId: 1}
	resp, err := client.GetBalance(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling GetBalance: %v", err)
	}
	log.Printf("Result: %v", resp)
}

func withdraw(client pb.PaymentsServiceClient) {
	req := &pb.WithdrawRequest{UserId: 1, Amount: 100}
	resp, err := client.Withdraw(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling Withdraw: %v", err)
	}
	log.Printf("Result: %v", resp)
}

func deposit(client pb.PaymentsServiceClient) {
	req := &pb.DepositRequest{UserId: 1, Amount: 100}
	resp, err := client.Deposit(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling Deposit: %v", err)
	}
	log.Printf("Result: %v", resp)
}

func getTransactionHistory(client pb.PaymentsServiceClient) {
	req := &pb.TransactionHistoryRequest{UserId: 1}
	stream, err := client.GetTransactionHistory(context.Background(), req)
	if err != nil {
		return
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error calling GetTransactionHistory: %v", err)
		}
		log.Printf("Result: %v", resp)
	}
}

func transferMultiple(client pb.PaymentsServiceClient) {
	reqs := []*pb.TransactRequest{{UserId: 1, Amount: 100, Type: pb.TransactionType_DEPOSIT}, {UserId: 1, Amount: 50, Type: pb.TransactionType_WITHDRAW}, {UserId: 1, Amount: 200, Type: pb.TransactionType_DEPOSIT}}
	stream, err := client.TransactMultiple(context.Background())
	if err != nil {
		return
	}
	for _, req := range reqs {
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("Error sending transaction: %v", err)
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error closing transaction stream")
	}
	log.Printf("Result: %v", resp)
}

func realtimeTransaction(client pb.PaymentsServiceClient) {
	reqs := []*pb.TransactRequest{{UserId: 1, Amount: 100, Type: pb.TransactionType_DEPOSIT}, {UserId: 1, Amount: 50, Type: pb.TransactionType_WITHDRAW}, {UserId: 3, Amount: 200, Type: pb.TransactionType_DEPOSIT}}
	stream, err := client.RealTimeTransfer(context.Background())
	if err != nil {
		return
	}
	for _, req := range reqs {
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("Error sending transaction: %v", err)
		}
		resp, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error closing transaction stream")
		}
		log.Printf("Result: %v", resp)
	}
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()
	client := pb.NewPaymentsServiceClient(conn)
	deposit(client)
}
