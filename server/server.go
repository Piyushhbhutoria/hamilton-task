package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Piyushhbhutoria/grpc-api/logger"
	"github.com/Piyushhbhutoria/grpc-api/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

type server struct {
	wallet.UnimplementedWalletServiceServer
}

func (s *server) CreateUserWallet(ctx context.Context, req *wallet.CreateUserWalletRequest) (*wallet.CreateUserWalletResponse, error) {
	// Create a new wallet for the user
	Wallet := &wallet.Wallet{
		UserId:   req.GetUserId(),
		Currency: req.GetCurrency(),
		Balance:  0,
	}

	// Store the wallet in the database
	walletID, err := CreateWallet(Wallet)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "creating wallet failed: %v", err)
	}

	Wallet.WalletId = walletID

	return &wallet.CreateUserWalletResponse{Wallet: Wallet}, nil
}

func (s *server) RecordTransaction(ctx context.Context, req *wallet.RecordTransactionRequest) (*wallet.RecordTransactionResponse, error) {
	// Retrieve the user's wallet from the database
	Wallet, err := GetWallet(req.GetUserId(), req.GetCurrency())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fetching wallet failed: %v", err)
	}

	// Add a new ledger to the wallet
	ledger := &wallet.Ledger{
		CreatedAt:       time.Now().String(),
		TransactionType: req.GetTransactionType(),
		Amount:          req.GetAmount(),
		Description:     req.GetDescription(),
		WalletId:        Wallet.GetWalletId(),
		Currency:        Wallet.GetCurrency(),
	}

	txValid := false

	// Update the wallet balance
	if ledger.GetTransactionType() == wallet.TransactionType_CREDIT.String() {
		Wallet.Balance += ledger.GetAmount()
		txValid = true
	} else if ledger.GetTransactionType() == wallet.TransactionType_DEBIT.String() {
		if Wallet.GetBalance() >= ledger.GetAmount() {
			Wallet.Balance -= ledger.GetAmount()
			txValid = true
		}
	}

	if txValid {
		// Store the updated wallet in the database
		err = AddTransaction(Wallet.GetBalance(), ledger)
		if err != nil {
			logger.LogMessage("error", "storing transaction failed: %v", err)
			return nil, status.Errorf(codes.Internal, "adding transaction failed: %v", err)
		}

		Wallet.Ledgers = append(Wallet.Ledgers, ledger)
		return &wallet.RecordTransactionResponse{Balance: fmt.Sprintf("updated balance: %.2f %s", Wallet.GetBalance(), req.GetCurrency())}, nil
	}

	logger.LogMessage("error", "transaction failed: %v", err)
	return nil, status.Errorf(codes.Internal, "invalid transaction: %v", err)
}

func (s *server) GetWalletSummary(ctx context.Context, req *wallet.GetWalletSummaryRequest) (*wallet.GetWalletSummaryResponse, error) {
	// Retrieve the user's wallet from the database
	Wallet, err := GetWallets(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fetching wallets failed: %v", err)
	}

	// Return the wallet balance and currency
	walletSummary := make([]*wallet.WalletSummary, len(Wallet))
	for i, w := range Wallet {
		walletSummary[i] = &wallet.WalletSummary{
			Balance:  w.GetBalance(),
			Currency: w.GetCurrency(),
		}
	}

	return &wallet.GetWalletSummaryResponse{WalletSummary: walletSummary}, nil
}

func (s *server) GetTransactionHistory(ctx context.Context, req *wallet.GetTransactionHistoryRequest) (*wallet.GetTransactionHistoryResponse, error) {
	// Return a paginated list of the user's ledgers
	offset := (req.GetPageNumber() - 1) * req.GetPageSize()
	limit := req.GetPageSize()

	// Retrieve the user's transactions from the database
	transactions, err := GetTransactions(req.GetUserId(), &limit, &offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fetching transactions failed: %v", err)
	}

	return &wallet.GetTransactionHistoryResponse{Ledgers: transactions}, nil
}

func Init() {
	// Create a listener on TCP port 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcLogger := grpc.UnaryInterceptor(logger.GRPCLogger)
	// Create a gRPC server object
	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     5 * time.Minute,
			MaxConnectionAge:      10 * time.Minute,
			MaxConnectionAgeGrace: 5 * time.Minute,
			Time:                  30 * time.Second,
			Timeout:               10 * time.Second,
		}),
		grpcLogger,
	)

	// Attach the Ping service to the server
	wallet.RegisterWalletServiceServer(s, &server{})

	// Serve gRPC Server
	logger.LogMessage("info", "Starting gRPC server on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
