package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "grpc-learn/proto"
	"log"
	"net"
)

type paymentServer struct {
	pb.UnimplementedPaymentsServiceServer
}

type Account struct {
	id      int32
	balance int32
}

var accounts = []Account{
	{id: 1, balance: 100},
	{id: 2, balance: 200},
	{id: 3, balance: 50},
}

func (s *paymentServer) GetBalance(ctx context.Context, req *pb.BalanceRequest) (*pb.BalanceResponse, error) {
	log.Printf("GetBalance request recieved: %v", req)
	for _, account := range accounts {
		if account.id == req.GetUserId() {
			return &pb.BalanceResponse{
				Amount: account.balance,
			}, nil
		}
	}
	return nil, fmt.Errorf("No account with user id %s", req.GetUserId())
}

func (s *paymentServer) Deposit(ctx context.Context, req *pb.DepositRequest) (*pb.DepositResponse, error) {
	log.Printf("Deposit request recieved: %v", req)
	for _, account := range accounts {
		if account.id == req.GetUserId() {
			account.balance += req.GetAmount()
			return &pb.DepositResponse{UserId: account.id, Balance: account.balance}, nil
		}
	}
	return nil, fmt.Errorf("No account with user id %s", req.GetUserId())
}

func (s *paymentServer) Withdraw(ctx context.Context, req *pb.WithdrawRequest) (*pb.WithdrawResponse, error) {
	log.Printf("Withdraw request recieved: %v", req)
	for _, account := range accounts {
		if account.id == req.GetUserId() {
			if account.balance >= req.GetAmount() {
				account.balance -= req.GetAmount()
				return &pb.WithdrawResponse{UserId: account.id, Balance: account.balance}, nil
			} else {
				return nil, fmt.Errorf("Account balance insufficient")
			}
		}
	}
	return nil, fmt.Errorf("No account with user id %s", req.GetUserId())
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Server started listening on port: %d", 50051)
	srv := grpc.NewServer()
	pb.RegisterPaymentsServiceServer(srv, &paymentServer{})
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
