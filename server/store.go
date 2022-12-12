package server

import (
	"encoding/json"

	"github.com/Piyushhbhutoria/grpc-api/logger"
	"github.com/Piyushhbhutoria/grpc-api/store"
	"github.com/Piyushhbhutoria/grpc-api/wallet"
)

func CreateWallet(Wallet *wallet.Wallet) (string, error) {
	db := store.GetSQL()

	var walletID string
	// Create a new wallet for the user
	logger.LogMessage("info", "creating new wallet for user: %s", Wallet.GetUserId())
	err := db.QueryRow(`INSERT INTO wallets (user_id, currency, balance) 
		VALUES ($1, $2, $3) RETURNING wallet_id`,
		Wallet.GetUserId(), Wallet.GetCurrency(), Wallet.GetBalance()).Scan(&walletID)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"wallets_user_id_currency_idx\"" {
			// if error in inserting fetch existing wallet
			logger.LogMessage("info", "wallet already exists for user: %s", Wallet.GetUserId())
			err = db.QueryRow(`SELECT wallet_id FROM wallets WHERE user_id = $1 AND currency = $2`,
				Wallet.GetUserId(), Wallet.GetCurrency()).Scan(&walletID)
			if err != nil {
				return walletID, err
			}
		}
		return walletID, err
	}

	return walletID, nil
}

func GetWallet(userID, currency string) (*wallet.Wallet, error) {
	db := store.GetSQL()
	// Retrieve the user's wallet from the database
	var Wallet wallet.Wallet
	logger.LogMessage("info", "retrieving wallet for user: %s", userID)
	var ledgerJSON json.RawMessage
	err := db.QueryRow(`SELECT wallet_id, user_id, currency, balance, ledger
		FROM wallet_view
		WHERE user_id = $1 AND currency = $2`,
		userID, currency).Scan(&Wallet.WalletId, &Wallet.UserId, &Wallet.Currency, &Wallet.Balance, &ledgerJSON)
	if err != nil {
		return nil, err
	}

	if ledgerJSON == nil {
		return &Wallet, nil
	}

	// Unmarshal the ledger JSON into the wallet
	err = json.Unmarshal(ledgerJSON, &Wallet.Ledgers)
	if err != nil {
		return nil, err
	}

	return &Wallet, nil
}

func GetWallets(userID string) ([]*wallet.Wallet, error) {
	db := store.GetSQL()
	// Retrieve the user's wallets from the database
	var Wallets []*wallet.Wallet
	logger.LogMessage("info", "retrieving all wallet for user: %s", userID)
	rows, err := db.Query(`SELECT wallet_id, user_id, currency, balance
		FROM wallets
		WHERE user_id = $1`,
		userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var Wallet wallet.Wallet
		err = rows.Scan(&Wallet.WalletId, &Wallet.UserId, &Wallet.Currency, &Wallet.Balance)
		if err != nil {
			return nil, err
		}

		Wallets = append(Wallets, &Wallet)
	}

	return Wallets, nil
}

func AddTransaction(balance float32, ledger *wallet.Ledger) error {
	db := store.GetSQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	logger.LogMessage("info", "adding transaction to wallet: %s", ledger.GetWalletId())
	// Add a new ledger transaction to the wallet
	_, err = tx.Exec(`INSERT INTO ledger (wallet_id, transaction_type, amount, description)
		VALUES ($1, $2, $3, $4)`,
		ledger.GetWalletId(), ledger.GetTransactionType(), ledger.GetAmount(), ledger.GetDescription())
	if err != nil {
		return err
	}

	// Update the user's wallet balance in the database
	_, err = tx.Exec(`UPDATE wallets SET balance = $1 WHERE wallet_id = $2`, balance, ledger.GetWalletId())
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func GetTransactions(userID string, limit, offset *int32) ([]*wallet.Ledger, error) {
	db := store.GetSQL()

	// Retrieve the user's wallets from the database
	logger.LogMessage("info", "retrieving transactions for user: %s", userID)
	var Ledgers []*wallet.Ledger
	rows, err := db.Query(`SELECT created_at, transaction_type, amount, description, wallet_id, currency
		FROM ledger_view
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		userID, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var Ledger wallet.Ledger
		err = rows.Scan(&Ledger.CreatedAt, &Ledger.TransactionType, &Ledger.Amount, &Ledger.Description, &Ledger.WalletId, &Ledger.Currency)
		if err != nil {
			return nil, err
		}

		Ledgers = append(Ledgers, &Ledger)
	}

	return Ledgers, nil
}
