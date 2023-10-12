package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "grpc-learn/proto"
	"io"
	"log"
	"net"
)

type Account struct {
	id      int32
	balance int32
}

type Transaction struct {
	transactionId   int32
	userId          int32
	amount          int32
	transactionType pb.TransactionType
}

var accounts = []Account{
	{id: 1, balance: 100},
	{id: 2, balance: 200},
	{id: 3, balance: 50},
}

var transactions []Transaction

var counter int32

type paymentServer struct {
	pb.UnimplementedPaymentsServiceServer
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

func deposit(userId int32, amount int32) (int32, error) {
	for _, account := range accounts {
		if account.id == userId {
			account.balance += amount
			transaction := Transaction{transactionId: counter, userId: userId, amount: amount, transactionType: pb.TransactionType_DEPOSIT}
			transactions = append(transactions, transaction)
			counter++
			return account.balance, nil
		}
	}
	return 0, fmt.Errorf("No account with user id %s", userId)
}

func withdraw(userId int32, amount int32) (int32, error) {
	for _, account := range accounts {
		if account.id == userId {
			if account.balance >= amount {
				account.balance -= amount
				transaction := Transaction{transactionId: counter, userId: userId, amount: amount, transactionType: pb.TransactionType_WITHDRAW}
				transactions = append(transactions, transaction)
				counter++
				return account.balance, nil
			} else {
				return 0, fmt.Errorf("Account balance insufficient")
			}
		}
	}
	return 0, fmt.Errorf("No account with user id %s", userId)
}

func transact(userId int32, amount int32, transactType pb.TransactionType) error {
	for _, account := range accounts {
		if account.id == userId {
			switch transactType {
			case pb.TransactionType_DEPOSIT:
				_, err := deposit(userId, amount)
				if err != nil {
					return err
				}
			case pb.TransactionType_WITHDRAW:
				_, err := withdraw(userId, amount)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *paymentServer) Deposit(ctx context.Context, req *pb.DepositRequest) (*pb.DepositResponse, error) {
	log.Printf("Deposit request recieved: %v", req)
	balance, err := deposit(req.GetUserId(), req.GetAmount())
	if err != nil {
		return nil, err
	}
	return &pb.DepositResponse{UserId: req.GetUserId(), Balance: balance}, nil
}

func (p *paymentServer) Withdraw(ctx context.Context, req *pb.WithdrawRequest) (*pb.WithdrawResponse, error) {
	log.Printf("Withdraw request recieved: %v", req)
	balance, err := withdraw(req.GetUserId(), req.GetAmount())
	if err != nil {
		return nil, err
	}
	return &pb.WithdrawResponse{UserId: req.GetUserId(), Balance: balance}, nil
}

func (p *paymentServer) GetTransactionHistory(req *pb.TransactionHistoryRequest, server pb.PaymentsService_GetTransactionHistoryServer) error {
	log.Printf("Get Transaction History request recieved: %v", req)
	for _, transaction := range transactions {
		if req.GetUserId() == transaction.userId {
			response := &pb.TransactionHistoryResponse{TransactionId: transaction.transactionId, UserId: transaction.userId, Type: transaction.transactionType, Amount: transaction.amount}
			if err := server.Send(response); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *paymentServer) TransactMultiple(server pb.PaymentsService_TransactMultipleServer) error {
	log.Printf("Transact Multiple request recieved")
	for {
		req, err := server.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err := transact(req.GetUserId(), req.GetAmount(), req.GetType()); err != nil {
			return err
		}
	}
	if err := server.SendAndClose(&pb.TransactResponse{Success: true}); err != nil {
		return err
	}
	return nil
}

func (p *paymentServer) RealTimeTransfer(server pb.PaymentsService_RealTimeTransferServer) error {
	log.Printf("Reatime transfer request recieved")
	for {
		req, err := server.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err := transact(req.GetUserId(), req.GetAmount(), req.GetType()); err != nil {
			return err
		}
		if err := server.Send(&pb.TransactResponse{Success: true}); err != nil {
			return err
		}
	}
	return nil
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
