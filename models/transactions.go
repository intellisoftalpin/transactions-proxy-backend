package models

import "github.com/intellisoftalpin/transactions-proxy-backend/models/cwalletapi"

// Struct to create user`s single transaction
type CreateTransactionRequest struct {
	Type string          `json:"type"`
	Data TransactionData `json:"data"`
}

// Struct to store user`s new single transaction status
type TransactionStatus struct {
	Status string `json:"status"`
}

// Type to store all user`s transactions
type AllTransactions struct {
	Transactions []Transaction `json:"transactions"`
}

// Struct to store one transaction out of all owned by the user
type Transaction struct {
	ID        int    `json:"id"`
	Hash      string `json:"hash"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	TransactionData
	DecodedTx cwalletapi.Transaction `json:"decodedTx"`
}

type TransactionData struct {
	AddressTo      string `json:"addressTo"`
	TransferAmount string `json:"transferAmount"` // Quantity of lovelaces
	AssetAmount    string `json:"assetAmount"`    // Quantity of assets
	AssetDecimals  string `json:"decimals"`
	PolicyID       string `json:"policyId"`
	AssetID        string `json:"assetId"`
	CBOR           string `json:"cbor"`
}

type ActiveTransaction struct {
	ID       uint64 `json:"id"`
	UserID   uint64 `json:"userId"`
	TxID     uint64 `json:"txId"`
	Step     string `json:"step"`
	Attempts int    `json:"attempts"`
}

type OngoingTransaction struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"userId"`
	TxID      uint64 `json:"txId"`
	Step      string `json:"step"`
	Attempts  int    `json:"attempts"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
