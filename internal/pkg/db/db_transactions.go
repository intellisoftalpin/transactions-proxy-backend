package db

import (
	"database/sql"

	models "github.com/intellisoftalpin/transactions-proxy-backend/models"
)

type TransactionsDB struct {
	db *sql.DB
}

func NewTransactionsDB(db *sql.DB) (transactionsDB *TransactionsDB) {
	return &TransactionsDB{db: db}
}

// ################################################################################
// Get user`s single transaction
// Select query to get user`s single transaction
const querySelectSingleTransaction = `
	SELECT 	tx_type,
			tx_status,
			tx_hash,
			tx_cbor,
			addr_to,
			transfer_amount,
			asset_amount,
			asset_decimals,
			policy_id,
			asset_id,
			created_at,
			updated_at
	FROM transactions
	WHERE id = $1 AND user_id = $2
`

// SelectTransaction - function to select user`s single transaction if exists
func (t *TransactionsDB) SelectTransaction(id uint64, userID uint64) (tx *models.SingleTransactionResponse, err error) {
	tx = &models.SingleTransactionResponse{}
	row := t.db.QueryRow(querySelectSingleTransaction, id, userID)
	err = row.Scan(
		&tx.Type,
		&tx.Status,
		&tx.Hash,
		&tx.CBOR,
		&tx.AddressTo,
		&tx.TransferAmount,
		&tx.AssetAmount,
		&tx.AssetDecimals,
		&tx.PolicyID,
		&tx.AssetID,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

// ################################################################################
// Get all user`s transactions
// Select query to get user`s single transaction
const querySelectAllTransactions = `
	SELECT 	id, 
			tx_type, 
			tx_status, 
			tx_hash,
			tx_cbor,
			addr_to,
			transfer_amount,
			asset_amount,
			asset_decimals,
			policy_id,
			asset_id,
			created_at, 
			updated_at
	FROM transactions
	WHERE user_id = $1
`

// SelectAllTransactions - function to select all user`s transactions
func (t *TransactionsDB) SelectAllTransactions(userID uint64) (allTransactions models.AllTransactions, err error) {
	rows, err := t.db.Query(querySelectAllTransactions, userID)
	if err != nil {
		return allTransactions, err
	}
	defer rows.Close()

	transactions := make([]models.Transaction, 0)

	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.Type,
			&transaction.Status,
			&transaction.Hash,
			&transaction.CBOR,
			&transaction.AddressTo,
			&transaction.TransferAmount,
			&transaction.AssetAmount,
			&transaction.AssetDecimals,
			&transaction.PolicyID,
			&transaction.AssetID,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return allTransactions, err
		}

		if transaction.Type == "buy" {
			transaction.AssetAmount = ""
		}

		transactions = append(transactions, transaction)
	}

	allTransactions.Transactions = transactions

	return allTransactions, nil
}

// ################################################################################
// Save user`s transaction data from JSON
// Insert query to save transaction data from JSON
const queryInsertTransactionData = `
	INSERT INTO transactions (
		id,
		user_id,
		tx_type,
		tx_status,
		tx_hash,
		tx_cbor,
		addr_to,
		transfer_amount,
		asset_amount,
		asset_decimals,
		policy_id,
		asset_id,
		created_at,
		updated_at
	)
	VALUES (
		DEFAULT,
		$1,
		$2,
		DEFAULT,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		CURRENT_TIMESTAMP,
		CURRENT_TIMESTAMP
	)
	RETURNING id
`

// InsertTransactionData - function to insert user`s new transaction data from JSON
func (t *TransactionsDB) InsertTransactionData(userID uint64, txType, txHash string, txData models.TransactionData) (transactionID uint64, err error) {
	row := t.db.QueryRow(queryInsertTransactionData,
		userID,
		txType,
		txHash,
		txData.CBOR,
		txData.AddressTo,
		txData.TransferAmount,
		txData.AssetAmount,
		txData.AssetDecimals,
		txData.PolicyID,
		txData.AssetID,
	)

	if err = row.Scan(&transactionID); err != nil {
		return transactionID, err
	}

	return transactionID, nil
}

// Update query to save new transaction data from JSON
const queryUpdateTransactionData = `
	UPDATE transactions
	SET tx_type = $3,
		tx_cbor = $4,
		updated_at = CURRENT_TIMESTAMP,
		addr_to = $5,
		transfer_amount = $6,
		policy_id = $7,
		asset_id = $8
	WHERE id = $1 AND user_id = $2
`

// UpdateTransactionData - function to update user`s transaction data from JSON
func (t *TransactionsDB) UpdateTransactionData(txID uint64, userID uint64, txType string, txData models.TransactionData) (err error) {
	query, err := t.db.Prepare(queryUpdateTransactionData)
	if err != nil {
		return err
	}
	_, err = query.Exec(txID, userID, txType,
		txData.CBOR,
		txData.AddressTo,
		txData.TransferAmount,
		txData.PolicyID,
		txData.AssetID,
	)
	if err != nil {
		return err
	}

	return nil
}

// Update query to save new transaction data from JSON
const queryUpdateTransactionHash = `
	UPDATE transactions
	SET tx_hash = $3,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $1 AND user_id = $2
`

// UpdateTransactionHash - function to update user`s transaction hash
func (t *TransactionsDB) UpdateTransactionHash(txID uint64, userID uint64, txHash string) (err error) {
	query, err := t.db.Prepare(queryUpdateTransactionHash)
	if err != nil {
		return err
	}

	_, err = query.Exec(txID, userID, txHash)
	if err != nil {
		return err
	}

	return nil
}

// ################################################################################
// Get user`s single transaction status
// Select query to get user`s single transaction status
const querySelectSingleTransactionStatus = `
	SELECT tx_status
	FROM transactions
	WHERE id = $1 AND user_id = $2
`

// SelectSingleTransactionStatus - function to select user`s single transaction status
func (t *TransactionsDB) SelectSingleTransactionStatus(txID uint64, userID uint64) (singleTransactionStatusResponse *models.SingleTransactionStatusResponse, err error) {
	singleTransactionStatusResponse = &models.SingleTransactionStatusResponse{}
	row := t.db.QueryRow(querySelectSingleTransactionStatus, txID, userID)

	err = row.Scan(&singleTransactionStatusResponse.TransactionStatus)
	if err != nil {
		return singleTransactionStatusResponse, err
	}

	return singleTransactionStatusResponse, nil
}

// ################################################################################
// Change user`s single transaction status with data from JSON
// Update query to change user`s single transaction status
const queryUpdateTransactionStatus = `
	UPDATE transactions
	SET tx_status = $3
	WHERE id = $1 AND user_id = $2
`

// UpdateTransactionStatus - function to update user`s single transaction status
func (t *TransactionsDB) UpdateTransactionStatus(txID uint64, userID uint64, txStatus string) (err error) {
	query, err := t.db.Prepare(queryUpdateTransactionStatus)
	if err != nil {
		return err
	}

	if _, err = query.Exec(txID, userID, txStatus); err != nil {
		return err
	}

	return nil
}

// ################################################################################
// Delete user`s single transaction
// Query to delete user`s single transaction
const queryDeleteSingleTransaction = `
	DELETE FROM transactions
	WHERE id = $1 AND user_id = $2
`

// DeleteSingleTransaction - function to delete user`s single transaction
func (t *TransactionsDB) DeleteSingleTransaction(txID uint64, userID uint64) (err error) {
	query, err := t.db.Prepare(queryDeleteSingleTransaction)
	if err != nil {
		return err
	}

	if _, err = query.Exec(txID, userID); err != nil {
		return err
	}

	return nil
}

// ---------------------------------------------------------------------------------

const queryInsertActiveTransaction = `
	INSERT INTO active_transactions (
		user_id,
		tx_id,
		step)
	VALUES (
		$1,
		$2,
		$3
	)
`

func (t *TransactionsDB) InsertActiveTransaction(userID uint64, txID uint64, step string) (err error) {
	query, err := t.db.Prepare(queryInsertActiveTransaction)
	if err != nil {
		return err
	}

	if _, err = query.Exec(userID, txID, step); err != nil {
		return err
	}

	return nil
}

const querySelectAllActiveTransactions = `
	SELECT id, user_id, tx_id, step, attempts
	FROM active_transactions
`

func (t *TransactionsDB) GetAllActiveTransactions() (activeTransactions []models.ActiveTransaction, err error) {
	rows, err := t.db.Query(querySelectAllActiveTransactions)
	if err != nil {
		return activeTransactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var activeTransaction models.ActiveTransaction
		err := rows.Scan(
			&activeTransaction.ID,
			&activeTransaction.UserID,
			&activeTransaction.TxID,
			&activeTransaction.Step,
			&activeTransaction.Attempts,
		)
		if err != nil {
			return activeTransactions, err
		}

		activeTransactions = append(activeTransactions, activeTransaction)
	}

	return activeTransactions, nil
}

const querySelectActiveTransactionsByStep = `
	SELECT id, user_id, tx_id, step, attempts
	FROM active_transactions
	WHERE step = $1
`

func (t *TransactionsDB) GetActiveTransactionsByStep(step string) (activeTransactions []models.ActiveTransaction, err error) {
	rows, err := t.db.Query(querySelectActiveTransactionsByStep, step)
	if err != nil {
		return activeTransactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var activeTransaction models.ActiveTransaction
		err := rows.Scan(
			&activeTransaction.ID,
			&activeTransaction.UserID,
			&activeTransaction.TxID,
			&activeTransaction.Step,
			&activeTransaction.Attempts,
		)
		if err != nil {
			return activeTransactions, err
		}

		activeTransactions = append(activeTransactions, activeTransaction)
	}

	return activeTransactions, nil
}

const queryDeleteActiveTransactionByID = `
	DELETE FROM active_transactions
	WHERE id = $1
`

func (t *TransactionsDB) DeleteActiveTransactionByID(id uint64) (err error) {
	query, err := t.db.Prepare(queryDeleteActiveTransactionByID)
	if err != nil {
		return err
	}

	if _, err = query.Exec(id); err != nil {
		return err
	}

	return nil
}

const queryUpdateActiveTransactionStep = `
	UPDATE active_transactions
	SET step = $2,
		attempts = 0
	WHERE id = $1
`

func (t *TransactionsDB) UpdateActiveTransactionStep(id uint64, step string) (err error) {
	query, err := t.db.Prepare(queryUpdateActiveTransactionStep)
	if err != nil {
		return err
	}

	if _, err = query.Exec(id, step); err != nil {
		return err
	}

	return nil
}

const queryUpdateActiveTransactionAttempts = `
	UPDATE active_transactions
	SET attempts = $2
	WHERE id = $1
`

func (t *TransactionsDB) UpdateActiveTransactionAttempts(id uint64, attempts int) (err error) {
	query, err := t.db.Prepare(queryUpdateActiveTransactionAttempts)
	if err != nil {
		return err
	}

	if _, err = query.Exec(id, attempts); err != nil {
		return err
	}

	return nil
}

const querySelectOngoingTransactions = `
	SELECT id, user_id, tx_id, step, attempts, created_at, updated_at
	FROM active_transactions
	WHERE user_id = $1
`

func (t *TransactionsDB) GetOngoingTransactions(userID uint64) (ongoingTransactions []models.OngoingTransaction, err error) {
	rows, err := t.db.Query(querySelectOngoingTransactions, userID)
	if err != nil {
		return ongoingTransactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var ongoingTransaction models.OngoingTransaction
		err := rows.Scan(
			&ongoingTransaction.ID,
			&ongoingTransaction.UserID,
			&ongoingTransaction.TxID,
			&ongoingTransaction.Step,
			&ongoingTransaction.Attempts,
			&ongoingTransaction.CreatedAt,
			&ongoingTransaction.UpdatedAt,
		)
		if err != nil {
			return ongoingTransactions, err
		}

		ongoingTransactions = append(ongoingTransactions, ongoingTransaction)
	}

	return ongoingTransactions, nil
}
