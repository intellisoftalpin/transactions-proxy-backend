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

type PoolsAPI struct {
	poolsRepo *repository.PoolsRepo

	WalletClient walletPB.WalletClient
}

func NewPoolsAPI(config *models.Config, walletClient walletPB.WalletClient) *PoolsAPI {
	return &PoolsAPI{
		WalletClient: walletClient,

		poolsRepo: repository.NewPoolsRepo(config, walletClient),
	}
}

// ################################################################################
// GetAllPools - function to get all pools
// GetAllPools godoc
// @Summary Get All Pools.
// @Description Get all pools.
// @Tags pools
// @Produce json
// @Success 200 {object} models.Pools
// @Router /api/v1/pools [GET]
func (api *PoolsAPI) GetAllPools(c echo.Context) error {
	pools, err := api.poolsRepo.GetAllPools()
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, pools)
}

// ################################################################################
// DelegateToPool - function to delegate to pool
// DelegateToPool godoc
// @Summary Delegate to pool.
// @Description Delegate to pool.
// @Tags pools
// @Accept json
// @Produce json
// @Param DelegateToPoolRequest body models.DelegateToPoolRequest true "DelegateToPoolRequest"
// @Success 200 {object} string
// @Router /api/v1/pools/delegate [POST]
func (api *PoolsAPI) DelegateToPool(c echo.Context) error {
	defer c.Request().Body.Close()

	// Get transaction data from JSON to save into database
	var delegateToPool models.DelegateToPoolRequest

	err := json.NewDecoder(c.Request().Body).Decode(&delegateToPool)
	if err != nil {
		fmt.Printf("Failed reading the request body: %s", err)
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	_, err = api.poolsRepo.DelegateToPool(delegateToPool)
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return nil
}
