package server

import (
	"context"
	"math/rand"
	"testing"

	"github.com/Piyushhbhutoria/grpc-api/store"
	"github.com/Piyushhbhutoria/grpc-api/wallet"
	"github.com/stretchr/testify/require"
)

var userID string

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSring(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateUser() error {
	db := store.GetSQL()
	// create user
	username := randSring(7)
	var UserID string
	err := db.QueryRow("INSERT INTO users (username) VALUES ($1) RETURNING user_id", username).Scan(&UserID)
	if err != nil {
		return err
	}
	userID = UserID
	return nil
}

func TestCreateUserWallet(t *testing.T) {
	testClient := wallet.NewWalletServiceClient(cc)
	// Create wallet for user
	req := wallet.CreateUserWalletRequest{
		UserId:   userID,
		Currency: "ETH",
	}

	result, err := testClient.CreateUserWallet(context.Background(), &req)
	expected := wallet.Wallet{
		WalletId: result.GetWallet().GetWalletId(),
		UserId:   userID,
		Currency: "ETH",
		Balance:  0,
		Ledgers:  nil,
	}
	require.Nil(t, err)
	require.Equal(t, expected, *result.GetWallet())

	// Create new wallet for user
	req = wallet.CreateUserWalletRequest{
		UserId:   userID,
		Currency: "BTC",
	}

	result, err = testClient.CreateUserWallet(context.Background(), &req)
	expected = wallet.Wallet{
		WalletId: result.GetWallet().GetWalletId(),
		UserId:   userID,
		Currency: "BTC",
		Balance:  0,
		Ledgers:  nil,
	}
	require.Nil(t, err)
	require.Equal(t, expected, *result.GetWallet())

	// Create duplicate wallet for user
	req = wallet.CreateUserWalletRequest{
		UserId:   userID,
		Currency: "BTC",
	}

	result, err = testClient.CreateUserWallet(context.Background(), &req)
	expected = wallet.Wallet{
		WalletId: result.GetWallet().GetWalletId(),
		UserId:   userID,
		Currency: "BTC",
		Balance:  0,
		Ledgers:  nil,
	}
	require.Nil(t, err)
	require.Equal(t, expected, *result.GetWallet())
}

func TestRecordTransaction(t *testing.T) {
	testClient := wallet.NewWalletServiceClient(cc)
	// Add money to wallet
	req := wallet.RecordTransactionRequest{
		UserId:          userID,
		TransactionType: "CREDIT",
		Currency:        "ETH",
		Amount:          1000,
		Description:     "Add money to wallet",
	}

	result, err := testClient.RecordTransaction(context.Background(), &req)
	expected := "updated balance: 1000.00 ETH"
	require.Nil(t, err)
	require.Equal(t, expected, result.GetBalance())

	// Withdraw money from wallet
	req = wallet.RecordTransactionRequest{
		UserId:          userID,
		TransactionType: "DEBIT",
		Currency:        "ETH",
		Amount:          100,
		Description:     "withdraw money from wallet",
	}

	// Invalid transaction
	result, err = testClient.RecordTransaction(context.Background(), &req)
	expected = "updated balance: 900.00 ETH"
	require.Nil(t, err)
	require.Equal(t, expected, result.GetBalance())

	req = wallet.RecordTransactionRequest{
		UserId:          userID,
		TransactionType: "DEBIT",
		Currency:        "ETH",
		Amount:          1000,
		Description:     "invalid transaction",
	}

	_, err = testClient.RecordTransaction(context.Background(), &req)
	require.NotNil(t, err)
}

func TestGetWalletSummary(t *testing.T) {
	testClient := wallet.NewWalletServiceClient(cc)
	req := wallet.GetWalletSummaryRequest{
		UserId: userID,
	}

	result, err := testClient.GetWalletSummary(context.Background(), &req)
	expected := []*wallet.WalletSummary{
		{
			Balance:  900.00,
			Currency: "ETH",
		},
		{
			Balance:  0.00,
			Currency: "BTC",
		},
	}
	require.Nil(t, err)
	require.Equal(t, expected, result.GetWalletSummary())

	// user not found
	req = wallet.GetWalletSummaryRequest{
		UserId: "test",
	}

	_, err = testClient.GetWalletSummary(context.Background(), &req)
	require.NotNil(t, err)
}

func TestGetTransactionHistory(t *testing.T) {
	testClient := wallet.NewWalletServiceClient(cc)
	req := wallet.GetTransactionHistoryRequest{
		UserId:     userID,
		PageSize:   10,
		PageNumber: 1,
	}

	result, err := testClient.GetTransactionHistory(context.Background(), &req)
	require.Nil(t, err)
	require.Equal(t, 2, len(result.GetLedgers()))
}
