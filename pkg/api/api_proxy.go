package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	walletPB "github.com/intellisoftalpin/proto/proto-gen/wallet"
	"github.com/intellisoftalpin/transactions-proxy-backend/consts"
	"github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/repository"
	utils "github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/utils"
	models "github.com/intellisoftalpin/transactions-proxy-backend/models"
	"github.com/labstack/echo/v4"
)

type ProxyAPI struct {
	proxyRepo *repository.ProxyRepo

	WalletClient walletPB.WalletClient
}

func NewProxyAPI(config *models.Config, walletClient walletPB.WalletClient) *ProxyAPI {
	return &ProxyAPI{
		WalletClient: walletClient,

		// poolsRepo: repository.NewPoolsRepo(config, walletClient),
	}
}

// ################################################################################
// SubmitExternalTransaction - function to submit external transaction
// SubmitExternalTransaction godoc
// @Summary Submit external transaction.
// @Description Submit external transaction.
// @Tags proxy
// @Accept json
// @Produce json
// @Param SubmitExternalTransactionRequest body models.SubmitExternalTransactionRequest true "SubmitExternalTransactionRequest"
// @Success 200 {object} models.SubmitExternalTransactionResponse
// @Router /api/v1/proxy/transactions [POST]
func (api *ProxyAPI) SubmitExternalTransaction(c echo.Context) error {
	defer c.Request().Body.Close()

	// Get transaction data from JSON to save into database
	var request models.SubmitExternalTransactionRequest

	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		fmt.Printf("Failed reading the request body: %s", err)
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusBadRequest)
	}

	txHash, err := api.proxyRepo.SubmitExternalTransaction(request)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, models.SubmitExternalTransactionResponse{TxHash: txHash})
}
