package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	walletPB "github.com/intellisoftalpin/proto/proto-gen/wallet"

	"github.com/intellisoftalpin/transactions-proxy-backend/consts"
	tlpsdb "github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/db"
	models "github.com/intellisoftalpin/transactions-proxy-backend/models"
	"github.com/intellisoftalpin/transactions-proxy-backend/models/cwalletapi"
)

type TransactionsRepo struct {
	// TransactionsDB tlpsdb.TransactionsDBInterface
	TransactionsDB *tlpsdb.TransactionsDB

	UsersSessions *models.Sessions
	WalletClient  walletPB.WalletClient

	// transactionsToSubmit       *TransactionsQueue
	// transactionsToUpdateStatus *TransactionsQueue
	// transactionsToSend         *TransactionsQueue
}

func NewTransactionsRepo(db *sql.DB, config *models.Config, sessions *models.Sessions, walletClient walletPB.WalletClient) (transactionsRepo *TransactionsRepo) {
	t := &TransactionsRepo{
		TransactionsDB: tlpsdb.NewTransactionsDB(db),
		UsersSessions:  sessions,
		WalletClient:   walletClient,

		// transactionsToSubmit:       NewTransactionsQueue(),
		// transactionsToUpdateStatus: NewTransactionsQueue(),
		// transactionsToSend:         NewTransactionsQueue(),
	}

	go t.submitTransactions()
	go t.updateTransactionStatus()
	go t.createTransactions()

	return t
}

// ################################################################################
// GetAllTransactions - function to get all user`s transactions
func (t *TransactionsRepo) GetAllTransactions(sessionID string) (allTransactionsResponse models.AllTransactions, err error) {
	userID, err := t.GetUserIDFromSessionID(sessionID)
	if err != nil {
		return allTransactionsResponse, err
	}

	allTransactionsResponse, err = t.TransactionsDB.SelectAllTransactions(userID)
	if err != nil {
		return allTransactionsResponse, err
	}

	ctx := context.Background()

	for i, tx := range allTransactionsResponse.Transactions {
		if tx.Hash != "" {

			resp, err := t.WalletClient.GetTransaction(ctx, &walletPB.GetTransactionRequest{
				TxHash:   tx.Hash,
				PolicyId: tx.PolicyID,
				AssetId:  tx.AssetID,
			})
			if err != nil {
				log.Println("GetTransaction TxID:", tx.ID, "UserID:", userID, "Error:", err)
				continue
			}

			var decodedTx cwalletapi.Transaction
			if err = json.Unmarshal(resp.RawTx, &decodedTx); err != nil {
				log.Println("Unmarshal TxID:", tx.ID, "UserID:", userID, "Error:", err)
				continue
			}

			allTransactionsResponse.Transactions[i].DecodedTx = decodedTx
		}
	}

	return allTransactionsResponse, nil
}

// ################################################################################
// GetSingleTransaction - function to get user`s single transaction
func (t *TransactionsRepo) GetSingleTransaction(txID uint64, sessionID string) (tx *models.SingleTransactionResponse, err error) {
	tx = &models.SingleTransactionResponse{}

	userID, err := t.GetUserIDFromSessionID(sessionID)
	if err != nil {
		return tx, err
	}

	if tx, err = t.TransactionsDB.SelectTransaction(txID, userID); err != nil {
		if err == sql.ErrNoRows {
			return tx, fmt.Errorf("transaction with id %d does not exist", txID)
		}

		return tx, err
	}

	if tx.Hash != "" {
		ctx := context.Background()

		resp, err := t.WalletClient.GetTransaction(ctx, &walletPB.GetTransactionRequest{
			TxHash:   tx.Hash,
			PolicyId: tx.PolicyID,
			AssetId:  tx.AssetID,
		})
		if err != nil {
			log.Println("GetTransaction TxID:", tx.ID, "UserID:", userID, "Error:", err)
			return tx, err
		}

		var decodedTx cwalletapi.Transaction
		if err = json.Unmarshal(resp.RawTx, &decodedTx); err != nil {
			log.Println("Unmarshal TxID:", tx.ID, "UserID:", userID, "Error:", err)
			return tx, err
		}

		tx.DecodedTx = decodedTx
	}

	return tx, nil
}

// ################################################################################
// GetSingleTransactionStatus - function to tgt user`s single transaction status
func (t *TransactionsRepo) GetSingleTransactionStatus(txID uint64, sessionID string) (singleTransactionStatusResponse *models.SingleTransactionStatusResponse, err error) {
	singleTransactionStatusResponse = &models.SingleTransactionStatusResponse{}

	userID, err := t.GetUserIDFromSessionID(sessionID)
	if err != nil {
		return singleTransactionStatusResponse, err
	}

	if singleTransactionStatusResponse, err = t.TransactionsDB.SelectSingleTransactionStatus(txID, userID); err != nil {
		if err == sql.ErrNoRows {
			return singleTransactionStatusResponse, fmt.Errorf("transaction with id %d does not exist", txID)
		}
		return singleTransactionStatusResponse, err
	}

	return singleTransactionStatusResponse, nil
}

// --------------------------------------------------------------------------------

// ################################################################################
// CreateTransaction - function to create user`s transaction data from JSON
func (t *TransactionsRepo) CreateTransaction(tx models.CreateTransactionRequest, sessionID string) (saveTransactionDataResponse *models.SaveTransactionDataResponse, err error) {
	saveTransactionDataResponse = &models.SaveTransactionDataResponse{}

	userID, err := t.GetUserIDFromSessionID(sessionID)
	if err != nil {
		return saveTransactionDataResponse, err
	}

	// Insert new transaction data
	saveTransactionDataResponse.TransactionID, err = t.TransactionsDB.InsertTransactionData(userID, tx.Type, "", tx.Data)
	if err != nil {
		return saveTransactionDataResponse, err
	}

	if tx.Type == "buy" {
		_, err = t.WalletClient.CheckTokenBalance(context.Background(), &walletPB.CheckTokenBalanceRequest{
			Tx:       tx.Data.CBOR,
			PolicyId: tx.Data.PolicyID,
			AssetId:  tx.Data.AssetID,
		})
		if err != nil {
			return saveTransactionDataResponse, err
		}

		_, err = t.UpdateTransactionStatus(saveTransactionDataResponse.TransactionID, models.TransactionStatus{
			Status: consts.CTransactionStatusPrepared}, sessionID)
		if err != nil {
			return saveTransactionDataResponse, err
		}
	}

	return saveTransactionDataResponse, nil
}

func (t *TransactionsRepo) UpdateTransaction(txID uint64, txType string, txData models.TransactionData, sessionID string) (singleTransactionResponse *models.SingleTransactionResponse, err error) {
	userID, err := t.GetUserIDFromSessionID(sessionID)
	if err != nil {
		return singleTransactionResponse, err
	}

	singleTxStatus, err := t.TransactionsDB.SelectSingleTransactionStatus(txID, userID)
	if err != nil {
		return singleTransactionResponse, err
	}

	if singleTxStatus.TransactionStatus != consts.CTransactionStatusDraft {
		return singleTransactionResponse, fmt.Errorf("unable to save transactions with status: %s, allowed statuses: draft, prepared", singleTxStatus.TransactionStatus)
	}

	if err = t.TransactionsDB.UpdateTransactionData(txID, userID, txType, txData); err != nil {
		return singleTransactionResponse, err
	}

	if singleTransactionResponse, err = t.TransactionsDB.SelectTransaction(txID, userID); err != nil {
		return singleTransactionResponse, err
	}

	return singleTransactionResponse, nil
}

// ################################################################################
// UpdateTransactionStatus - function to change user`s single transaction status with data from JSON
func (t *TransactionsRepo) UpdateTransactionStatus(txID uint64, txStatus models.TransactionStatus, sessionID string) (changeSingleTransactionStatusResponse *models.ChangeSingleTransactionStatusResponse, err error) {
	changeSingleTransactionStatusResponse = &models.ChangeSingleTransactionStatusResponse{}

	userID, err := t.GetUserIDFromSessionID(sessionID)
	if err != nil {
		return changeSingleTransactionStatusResponse, err
	}

	// Check transaction status. Propper changes draft -> prepared, prepared -> draft/submited, submited -> success/failure
	singleTransactionStatus, err := t.TransactionsDB.SelectSingleTransactionStatus(txID, userID)
	if err != nil {
		return changeSingleTransactionStatusResponse, err
	}

	if singleTransactionStatus.TransactionStatus == consts.CTransactionStatusDraft && txStatus.Status == consts.CTransactionStatusPrepared {
		if err = t.TransactionsDB.UpdateTransactionStatus(txID, userID, txStatus.Status); err != nil {
			return changeSingleTransactionStatusResponse, err
		}

		if err = t.TransactionsDB.InsertActiveTransaction(userID, txID, consts.CActiveTransactionStepSubmit); err != nil {
			return changeSingleTransactionStatusResponse, err
		}

		// t.transactionsToSubmit.Store(txID, userID, true)
		log.Println("Transaction TxID:", txID, "UserID:", userID, "Stored in transactionsToSubmit")

		changeSingleTransactionStatusResponse.Status = txStatus.Status
		// changeSingleTransactionStatusResponse.NewTransactionData = txStatus.NewTransactionData
		return changeSingleTransactionStatusResponse, nil
	}

	err = fmt.Errorf("unable to change transaction status from %s to %s", singleTransactionStatus.TransactionStatus, txStatus.Status)
	return changeSingleTransactionStatusResponse, err
}

// ################################################################################
// DeleteSingleTransaction - function to delete user`s single transaction
func (t *TransactionsRepo) DeleteSingleTransaction(txID uint64, sessionID string) (deleteSingleTransactionResponse *models.DeleteSingleTransactionResponse, err error) {
	deleteSingleTransactionResponse = &models.DeleteSingleTransactionResponse{}

	userID, err := t.GetUserIDFromSessionID(sessionID)
	if err != nil {
		deleteSingleTransactionResponse.Status = consts.CErrorStatus
		deleteSingleTransactionResponse.Message = consts.MessageDeleteTransactionError
		return deleteSingleTransactionResponse, err
	}

	// Check transaction status. Delete transaction only if status is not submited/success/failure
	singleTransactionStatus, err := t.TransactionsDB.SelectSingleTransactionStatus(txID, userID)
	if err != nil {
		deleteSingleTransactionResponse.Status = consts.CErrorStatus
		deleteSingleTransactionResponse.Message = consts.MessageDeleteTransactionError
		return deleteSingleTransactionResponse, err
	}

	if singleTransactionStatus.TransactionStatus == consts.CTransactionStatusDraft {
		if err = t.TransactionsDB.DeleteSingleTransaction(txID, userID); err != nil {
			log.Println("DeleteSingleTransaction TxID:", txID, "UserID:", userID, "Error:", err)
			deleteSingleTransactionResponse.Status = consts.CErrorStatus
			deleteSingleTransactionResponse.Message = consts.MessageDeleteTransactionError
			return deleteSingleTransactionResponse, err
		}

		deleteSingleTransactionResponse.Status = consts.CSuccessStatus
		deleteSingleTransactionResponse.Message = consts.MessageDeleteTransactionSuccess
		return deleteSingleTransactionResponse, nil
	}

	err = fmt.Errorf("could not delete transaction with status %s, could delete only if status is not submited/success/failure", singleTransactionStatus.TransactionStatus)
	deleteSingleTransactionResponse.Status = consts.CErrorStatus
	deleteSingleTransactionResponse.Message = consts.MessageDeleteTransactionError
	return deleteSingleTransactionResponse, err
}

// ################################################################################
// models.ActiveTransactionsResponse
// CheckActiveTransactions - function to get active user`s transactions
func (t *TransactionsRepo) CheckActiveTransactions(sessionID string) (activeTransactionsResponse models.ActiveTransactionsResponse, err error) {
	userID, err := t.GetUserIDFromSessionID(sessionID)
	if err != nil {
		return activeTransactionsResponse, err
	}

	ongoingTransactions, err := t.TransactionsDB.GetOngoingTransactions(userID)
	if err != nil {
		return activeTransactionsResponse, err
	}

	if len(ongoingTransactions) > 0 {
		activeTransactionsResponse.IsBusy = true
	}

	return activeTransactionsResponse, nil
}

// ------------------------------------------------------------------------------------------

func (t *TransactionsRepo) submitTransactions() {
	ticker := time.NewTicker(10 * time.Second)

	for {
		<-ticker.C

		activeTransactions, err := t.TransactionsDB.GetActiveTransactionsByStep(consts.CActiveTransactionStepSubmit)
		if err != nil {
			log.Println("GetAllActiveTransactions Error:", err)
			continue
		}

		for _, value := range activeTransactions {
			txID := value.TxID
			userID := value.UserID

			if value.Attempts > 5 {
				// 	t.transactionsToSubmit.Delete(txID)
				// 	log.Println("Transaction TxID:", txID, "UserID:", value.UserID, "Deleted from transactionsToSubmit")

				// Delete transaction from active_transactions
				if err = t.TransactionsDB.DeleteActiveTransactionByID(value.ID); err != nil {
					log.Println("DeleteActiveTransactionByID TxID:", txID, "UserID:", userID, "Error:", err)
					continue
				}

				// Update transaction status to "failed"
				if err := t.TransactionsDB.UpdateTransactionStatus(txID, userID, consts.CTransactionStatusFailed); err != nil {
					log.Println("UpdateTransactionStatus TxID:", txID, "UserID:", userID, "Error:", err)
					continue
				}

				continue
			}
			// t.transactionsToSubmit.Update(txID, value.Attempt-1)

			// Update transaction attempts in active_transactions
			if err = t.TransactionsDB.UpdateActiveTransactionAttempts(value.ID, value.Attempts+1); err != nil {
				log.Println("UpdateActiveTransactionAttempts TxID:", txID, "UserID:", userID, "Error:", err)
				continue
			}

			// check balance before submit transaction

			transaction, err := t.TransactionsDB.SelectTransaction(txID, userID)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Println(fmt.Errorf("transaction with id %d does not exist", txID))
					continue
				}
				log.Println("SelectTransaction TxID:", txID, "UserID:", userID, "Error:", err)
				continue
			}

			_, err = t.WalletClient.CheckTokenBalance(context.Background(), &walletPB.CheckTokenBalanceRequest{
				Tx:       transaction.CBOR,
				PolicyId: transaction.PolicyID,
				AssetId:  transaction.AssetID,
			})
			if err != nil {
				log.Println("CheckTokenBalance TxID:", txID, "UserID:", userID, "Error:", err)

				// Update transaction status to "failed"
				if err = t.TransactionsDB.UpdateTransactionStatus(txID, userID, consts.CTransactionStatusFailed); err != nil {
					log.Println("UpdateTransactionStatus TxID:", txID, "UserID:", userID, "Error:", err)
					continue
				}

				// Delete transaction from active_transactions
				if err = t.TransactionsDB.DeleteActiveTransactionByID(value.ID); err != nil {
					log.Println("DeleteActiveTransactionByID TxID:", txID, "UserID:", userID, "Error:", err)
					continue
				}

				// t.transactionsToSubmit.Delete(txID)
				log.Println("Transaction TxID:", txID, "UserID:", userID, "Deleted from transactionsToSubmit")
				continue
			}

			if err = t.submitTransaction(txID, userID); err != nil {
				log.Println("SubmitTransaction TxID:", txID, "UserID:", userID, "Error:", err)
				continue
			}

			// Update transaction status from active_transactions
			if err = t.TransactionsDB.UpdateActiveTransactionStep(value.ID, consts.CActiveTransactionStepUpdateStatus); err != nil {
				log.Println("UpdateActiveTransactionStep TxID:", txID, "UserID:", userID, "Error:", err)
				continue
			}

			// t.transactionsToUpdateStatus.Store(txID, userID, true)
			// log.Println("Transaction TxID:", txID, "UserID:", userID, "Stored in transactionsToUpdateStatus")

			// t.transactionsToSubmit.Delete(txID)
			// log.Println("Transaction TxID:", txID, "UserID:", userID, "Deleted from transactionsToSubmit")
		}
	}
}

func (t *TransactionsRepo) submitTransaction(txID uint64, userID uint64) (err error) {
	tx, err := t.TransactionsDB.SelectTransaction(txID, userID)
	if err != nil {
		return err
	}
	ctx := context.Background()

	if tx.Status == consts.CTransactionStatusPrepared {
		// Submit transaction
		resp, err := t.WalletClient.SubmitTransaction(ctx, &walletPB.SubmitTransactionRequest{Tx: tx.CBOR})
		if err != nil {
			return err
		}

		err = t.TransactionsDB.UpdateTransactionStatus(txID, userID, consts.CTransactionStatusSubmitted)
		if err != nil {
			return err
		}

		err = t.TransactionsDB.UpdateTransactionHash(txID, userID, resp.TxHash)
		if err != nil {
			return err
		}

		return nil
	}

	err = fmt.Errorf("unable to submit transaction with status %s", tx.Status)
	return err
}

func (t *TransactionsRepo) updateTransactionStatus() {
	ticker := time.NewTicker(10 * time.Second)

	for {
		<-ticker.C

		activeTransactions, err := t.TransactionsDB.GetActiveTransactionsByStep(consts.CActiveTransactionStepUpdateStatus)
		if err != nil {
			log.Println("GetAllActiveTransactions Error:", err)
			continue
		}

		activeTransactionsOnlyUpdateStatus, err := t.TransactionsDB.GetActiveTransactionsByStep(consts.CActiveTransactionStepOnlyUpdateStatus)
		if err != nil {
			log.Println("GetAllActiveTransactions Error:", err)
			continue
		}

		activeTransactions = append(activeTransactions, activeTransactionsOnlyUpdateStatus...)

		for _, value := range activeTransactions {
			txID := value.TxID
			userID := value.UserID

			status, err := t.checkTransactionStatus(txID, userID)
			if err != nil {
				log.Println("checkTransactionStatus TxID:", txID, "UserID:", userID, "Error:", err)
				continue
			}

			if err = t.TransactionsDB.UpdateTransactionStatus(txID, userID, status); err != nil {
				log.Println("updateTransactionStatus TxID:", txID, "UserID:", userID, "Error:", err)
				continue
			}

			if status == "in_ledger" || status == "expired" {
				if status == "in_ledger" && value.Step == consts.CActiveTransactionStepUpdateStatus {
					// t.transactionsToSend.Store(txID, userID, true)

					// Update transaction status from active_transactions
					if err = t.TransactionsDB.UpdateActiveTransactionStep(value.ID, consts.CActiveTransactionStepSend); err != nil {
						log.Println("UpdateActiveTransactionStep TxID:", txID, "UserID:", userID, "Error:", err)
						continue
					}

					log.Println("Transaction TxID:", txID, "UserID:", userID, "Stored in transactionsToSend")
				} else {
					// t.transactionsToUpdateStatus.Delete(txID)

					// Delete transaction from active_transactions
					if err = t.TransactionsDB.DeleteActiveTransactionByID(value.ID); err != nil {
						log.Println("DeleteActiveTransactionByID TxID:", txID, "UserID:", userID, "Error:", err)
						continue
					}

					log.Println("Transaction TxID:", txID, "UserID:", userID, "Status:", status, "Deleted from transactionsToUpdateStatus")
				}
			}
		}
	}
}

func (t *TransactionsRepo) checkTransactionStatus(txID uint64, userID uint64) (txStatus string, err error) {
	tx, err := t.TransactionsDB.SelectTransaction(txID, userID)
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	resp, err := t.WalletClient.GetTransaction(ctx, &walletPB.GetTransactionRequest{
		TxHash:   tx.Hash,
		PolicyId: tx.PolicyID,
		AssetId:  tx.AssetID,
	})
	if err != nil {
		return "", err
	}

	var walletTx WalletTransaction

	if err = json.Unmarshal(resp.RawTx, &walletTx); err != nil {
		return "", err
	}

	return walletTx.Status, nil
}

func (t *TransactionsRepo) createTransactions() {
	ticker := time.NewTicker(10 * time.Second)

	for {
		<-ticker.C

		activeTransactions, err := t.TransactionsDB.GetActiveTransactionsByStep(consts.CActiveTransactionStepSend)
		if err != nil {
			log.Println("GetAllActiveTransactions Error:", err)
			continue
		}

		for _, value := range activeTransactions {
			txID := value.TxID
			userID := value.UserID

			ctx := context.Background()

			tx, err := t.TransactionsDB.SelectTransaction(txID, userID)
			if err != nil {
				log.Println("SelectTransaction TxID:", txID, "UserID:", userID, "Error:", err)
				continue
			}

			resp, err := t.WalletClient.CreateTransaction(ctx, &walletPB.CreateTransactionRequest{
				Tx:       tx.CBOR,
				PolicyId: tx.PolicyID,
				AssetId:  tx.AssetID,
			})
			if err != nil {
				log.Println("CreateTransaction TxID:", txID, "UserID:", userID, "Error:", err)
				continue
			}

			log.Println("CreateTransaction TxID:", txID, "UserID:", userID, "Response:", resp)
			// t.transactionsToSend.Delete(txID)
			log.Println("Transaction TxID:", txID, "UserID:", userID, "Deleted from transactionsToSend")

			txID, err = t.TransactionsDB.InsertTransactionData(userID, "reverse", resp.TxHash, models.TransactionData{
				AddressTo:      resp.AddressTo,
				TransferAmount: resp.TransferAmount, // Quantity of lovelaces
				AssetAmount:    resp.AssetAmount,    // Quantity of tokens
				AssetDecimals:  resp.AssetDecimals,
				PolicyID:       tx.PolicyID,
				AssetID:        tx.AssetID,
				// CBOR:           resp.Tx,
			})
			if err != nil {
				log.Println("InsertTransactionData TxID:", txID, "UserID:", userID, "Error:", err)
				continue
			}

			// Delete transaction from active_transactions
			if err = t.TransactionsDB.DeleteActiveTransactionByID(value.ID); err != nil {
				log.Println("DeleteActiveTransactionByID TxID:", txID, "UserID:", userID, "Error:", err)
				continue
			}

			// Update transaction status from active_transactions
			if err = t.TransactionsDB.InsertActiveTransaction(userID, txID, consts.CActiveTransactionStepOnlyUpdateStatus); err != nil {
				log.Println("UpdateActiveTransactionStep TxID:", txID, "UserID:", userID, "Error:", err)
				continue
			}

			// t.transactionsToUpdateStatus.Store(txID, userID, false)

		}
	}
}

type WalletTransaction struct {
	Status string `json:"status"`
}

// ------------------------------------------------------------------------------------------

// GetUserIDFromSessionID - internal function for getting user ID from session ID
func (t *TransactionsRepo) GetUserIDFromSessionID(sessionID string) (userID uint64, err error) {
	// Get session from sessions map by Authorization Key (sessionID)
	session, exists := t.UsersSessions.SessionsMap[sessionID]

	// Check if session does not exist
	if !exists {
		return userID, fmt.Errorf("session with authorization key %s does not exist", sessionID)
	}

	// Check if session is expired
	if time.Now().After(session.ExpirationDateTime) {
		return userID, fmt.Errorf("session with authorization key %s is expired", sessionID)
	}

	// If session is not expired than get userID
	return session.UserID, nil
}

// ------------------------------------------------------------------------------------------

// type TransactionsQueue struct {
// 	mx *sync.RWMutex
// 	// txMap map[uint64]uint64 // key - transaction ID, value - user ID
// 	txMap map[uint64]TransactionsQueueItem // key - transaction ID
// }

// type TransactionsQueueItem struct {
// 	// TxID        uint64
// 	UserID      uint64
// 	CreateNewTx bool

// 	Attempt int
// }

// func NewTransactionsQueue() *TransactionsQueue {
// 	return &TransactionsQueue{
// 		mx:    &sync.RWMutex{},
// 		txMap: make(map[uint64]TransactionsQueueItem),
// 	}
// }

// func (tx *TransactionsQueue) Store(txID uint64, userID uint64, createNewTx bool) {
// 	tx.mx.Lock()
// 	defer tx.mx.Unlock()
// 	tx.txMap[txID] = TransactionsQueueItem{
// 		UserID:      userID,
// 		CreateNewTx: createNewTx,
// 		Attempt:     5,
// 	}
// }

// func (tx *TransactionsQueue) Update(txID uint64, attempt int) {
// 	tx.mx.Lock()
// 	defer tx.mx.Unlock()
// 	tx.txMap[txID] = TransactionsQueueItem{
// 		UserID:      tx.txMap[txID].UserID,
// 		CreateNewTx: tx.txMap[txID].CreateNewTx,
// 		Attempt:     attempt,
// 	}
// }

// func (tx *TransactionsQueue) Load(txID uint64) (value TransactionsQueueItem, exists bool) {
// 	tx.mx.RLock()
// 	defer tx.mx.RUnlock()
// 	value, exists = tx.txMap[txID]
// 	return value, exists
// }

// func (tx *TransactionsQueue) Delete(txID uint64) {
// 	tx.mx.Lock()
// 	defer tx.mx.Unlock()
// 	delete(tx.txMap, txID)
// }
