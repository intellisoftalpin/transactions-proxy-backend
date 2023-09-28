package api

import (
	"net/http"

	"github.com/intellisoftalpin/transactions-proxy-backend/consts"
	"github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/repository"
	utils "github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/utils"
	"github.com/labstack/echo/v4"
)

type PoolsAPI struct {
	poolsRepo *repository.PoolsRepo
}

func NewPoolsAPI() *PoolsAPI {
	return &PoolsAPI{
		poolsRepo: repository.NewPoolsRepo(),
	}
}

// ################################################################################
// GetAllPools - function to get all pools
// GetAllPools godoc
// @Summary Get All Pools.
// @Description Get all pools.
// @Tags pools
// @Produce json
// @Success 200 {object} []models.Pool
// @Router /api/v1/pools [GET]
func (api *PoolsAPI) GetAllPools(c echo.Context) error {
	pools, err := api.poolsRepo.GetAllPools()
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, pools)
}
