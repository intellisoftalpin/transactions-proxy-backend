package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	walletPB "github.com/intellisoftalpin/proto/proto-gen/wallet"
	"github.com/labstack/echo/v4"

	"github.com/intellisoftalpin/transactions-proxy-backend/consts"
	"github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/repository"
	utils "github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/utils"
	models "github.com/intellisoftalpin/transactions-proxy-backend/models"
)

type TransactionsAPI struct {
	transactionsRepo *repository.TransactionsRepo

	WalletClient walletPB.WalletClient
}

func NewTransactionsAPI(db *sql.DB, config *models.Config, sessions *models.Sessions, walletClient walletPB.WalletClient) *TransactionsAPI {
	return &TransactionsAPI{
		WalletClient: walletClient,

		transactionsRepo: repository.NewTransactionsRepo(db, config, sessions, walletClient),
	}
}

func (api *TransactionsAPI) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionID := c.Request().Header.Get("session_id")

		_, err := api.transactionsRepo.GetUserIDFromSessionID(sessionID)
		if err != nil {
			return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorInvalidSessionID, http.StatusUnauthorized)
		}

		return next(c)
	}
}

// ################################################################################
// GetAllTransactions - function to get all user`s transactions
// GetAllTransactions godoc
// @Summary Get All Transactions.
// @Description Get all user`s transaction.
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} models.AllTransactions
// @Router /api/v1/transactions [GET]
func (api *TransactionsAPI) GetAllTransactions(c echo.Context) error {
	sessionID := c.Request().Header.Get("session_id")
	allTransactionsResponse, err := api.transactionsRepo.GetAllTransactions(sessionID)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, allTransactionsResponse)
}

// ################################################################################
// GetSingleTransaction - function to get user`s single transaction
// GetSingleTransaction godoc
// @Summary Get Single Transaction.
// @Description Get user`s single transaction.
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} models.SingleTransactionResponse
// @Router /api/v1/transactions/:transaction_id [GET]
func (api *TransactionsAPI) GetSingleTransaction(c echo.Context) error {
	sessionID := c.Request().Header.Get("session_id")
	transactionID := c.Param("transaction_id")

	tempTransactionID, err := strconv.ParseUint(transactionID, 10, 64)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorInvalidTransactionID, http.StatusBadRequest)
	}

	singleTransactionResponse, err := api.transactionsRepo.GetSingleTransaction(tempTransactionID, sessionID)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, singleTransactionResponse)
}

// ################################################################################
// GetSingleTransactionStatus - function to get user`s single transaction status
// GetSingleTransactionStatus godoc
// @Summary Get Single Transaction Status.
// @Description Get user`s single transaction status.
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} models.SingleTransactionStatusResponse
// @Router /api/v1/transactions/:transaction_id/status [GET]
func (api *TransactionsAPI) GetSingleTransactionStatus(c echo.Context) error {
	transactionID := c.Param("transaction_id")

	sessionID := c.Request().Header.Get("session_id")

	tempTransactionID, err := strconv.ParseUint(transactionID, 10, 64)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorInvalidTransactionID, http.StatusBadRequest)
	}

	singleTransactionStatusResponse, err := api.transactionsRepo.GetSingleTransactionStatus(tempTransactionID, sessionID)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, singleTransactionStatusResponse)
}

// ################################################################################
// CreateTransaction - function to save user`s transaction data from JSON
// CreateTransaction godoc
// @Summary Save Transaction.
// @Description Save user`s transaction data from JSON.
// @Tags transactions
// @Accept json
// @Produce json
// @Param Transaction body models.CreateTransactionRequest true "Struct to create user`s single transaction"
// @Success 200 {object} models.SaveTransactionDataResponse
// @Router /api/v1/transactions [POST]
func (api *TransactionsAPI) CreateTransaction(c echo.Context) error {
	defer c.Request().Body.Close()

	sessionID := c.Request().Header.Get("session_id")

	// Get transaction data from JSON to save into database
	var cardanoTransaction models.CreateTransactionRequest

	err := json.NewDecoder(c.Request().Body).Decode(&cardanoTransaction)
	if err != nil {
		fmt.Printf("Failed reading the request body: %s", err)
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	saveTransactionDataResponse, err := api.transactionsRepo.CreateTransaction(cardanoTransaction, sessionID)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, saveTransactionDataResponse)
}

// ################################################################################
// UpdateSingleTransaction - function to update user`s single transaction data from JSON
// UpdateSingleTransaction godoc
// @Summary Update Transaction.
// @Description Update user`s single transaction data from JSON.
// @Tags transactions
// @Accept json
// @Produce json
// @Param Transaction body models.CreateTransactionRequest true "Struct to update user`s single transaction"
// @Success 200 {object} models.SaveTransactionDataResponse
// @Router /api/v1/transactions/:transaction_id [PUT]
func (api *TransactionsAPI) UpdateSingleTransaction(c echo.Context) error {
	transactionID := c.Param("transaction_id")

	sessionID := c.Request().Header.Get("session_id")

	tempTransactionID, err := strconv.ParseUint(transactionID, 10, 64)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorInvalidTransactionID, http.StatusBadRequest)
	}

	// Get transaction data from JSON to save into database
	var cardanoTransaction models.CreateTransactionRequest
	defer c.Request().Body.Close()

	err = json.NewDecoder(c.Request().Body).Decode(&cardanoTransaction)
	if err != nil {
		fmt.Printf("Failed reading the request body: %s", err)
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	saveTransactionDataResponse, err := api.transactionsRepo.UpdateTransaction(tempTransactionID, cardanoTransaction.Type, cardanoTransaction.Data, sessionID)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, saveTransactionDataResponse)
}

// ################################################################################
// ChangeSingleTransactionStatus - function to change user`s single transaction status with data from JSON
// ChangeSingleTransactionStatus godoc
// @Summary Change Single Transaction Status.
// @Description Change user`s single transaction status with data from JSON.
// @Tags transactions
// @Accept json
// @Produce json
// @Param TransactionNewStatus body models.TransactionStatus true "Struct to store user`s new single transaction status"
// @Success 200 {object} models.ChangeSingleTransactionStatusResponse
// @Router /api/v1/transactions/:transaction_id/status [PUT]
func (api *TransactionsAPI) ChangeSingleTransactionStatus(c echo.Context) error {
	transactionID := c.Param("transaction_id")

	req := c.Request()
	sessionID := req.Header.Get("session_id")
	body := req.Body

	tempTransactionID, err := strconv.ParseUint(transactionID, 10, 64)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorInvalidTransactionID, http.StatusBadRequest)
	}

	// Get transaction status data from JSON to update record in database table
	var singleTransactionNewStatus models.TransactionStatus
	defer body.Close()
	err = json.NewDecoder(body).Decode(&singleTransactionNewStatus)
	if err != nil {
		fmt.Printf("Failed reading the request body: %s", err)
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	changeSingleTransactionStatusResponse, err := api.transactionsRepo.UpdateTransactionStatus(tempTransactionID, singleTransactionNewStatus, sessionID)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, changeSingleTransactionStatusResponse)
}

// ################################################################################
// DeleteSingleTransaction - function to delete user`s single transaction
// DeleteSingleTransaction godoc
// @Summary Delete Single Transaction.
// @Description Delete user`s single transaction.
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} models.DeleteSingleTransactionResponse
// @Router /api/v1/transactions/:transaction_id [DELETE]
func (api *TransactionsAPI) DeleteSingleTransaction(c echo.Context) error {
	transactionID := c.Param("transaction_id")

	sessionID := c.Request().Header.Get("session_id")

	tempTransactionID, err := strconv.ParseUint(transactionID, 10, 64)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorInvalidTransactionID, http.StatusBadRequest)
	}

	deleteSingleTransactionResponse, err := api.transactionsRepo.DeleteSingleTransaction(tempTransactionID, sessionID)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, deleteSingleTransactionResponse)
}

// ################################################################################
// CheckActiveTransactions - function to check user`s active transactions
// CheckActiveTransactions godoc
// @Summary Check Active Transactions.
// @Description Check user`s active transactions.
// @Tags transactions
// @Produce json
// @Success 200 {object} models.ActiveTransactionsResponse
// @Router /api/v1/transactions/active [GET]
func (api *TransactionsAPI) CheckActiveTransactions(c echo.Context) error {
	sessionID := c.Request().Header.Get("session_id")

	ongoingTransactionsResponse, err := api.transactionsRepo.CheckActiveTransactions(sessionID)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, ongoingTransactionsResponse)
}
