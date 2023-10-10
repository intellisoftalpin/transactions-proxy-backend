package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	walletPB "github.com/intellisoftalpin/proto/proto-gen/wallet"
	"github.com/intellisoftalpin/transactions-proxy-backend/consts"
	utils "github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/utils"
	models "github.com/intellisoftalpin/transactions-proxy-backend/models"
	"github.com/labstack/echo/v4"
)

type NetworkAPI struct {
	WalletClient walletPB.WalletClient

	walletNetworkReady bool
	walletNetworkInfo  models.NetworkInfoResponse
}

func NewNetworkAPI(walletClient walletPB.WalletClient) *NetworkAPI {
	n := &NetworkAPI{
		WalletClient: walletClient,
	}

	go func() {
		ctx := context.Background()

		timer := time.NewTicker(5 * time.Second)

		for range timer.C {
			nInfo, err := n.WalletClient.GetWalletNetworkInfo(ctx, &walletPB.Empty{})
			if err != nil {
				n.walletNetworkReady = false
			} else {

				if nInfo.SyncProgress.Status == consts.CSyncProgressStatusReady {
					n.walletNetworkReady = true
				}

				n.walletNetworkInfo = models.ToNetworkInfoResponse(nInfo)
			}

		}
	}()

	return n
}

func (api *NetworkAPI) MiddlewareNetworkReady(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if !api.walletNetworkReady {
			msg := "Wallet network is not ready. Status: " + api.walletNetworkInfo.SyncProgress.Status

			if api.walletNetworkInfo.SyncProgress.Status == "syncing" {
				msg += ". Progress: " + strconv.FormatUint(api.walletNetworkInfo.SyncProgress.Progress.Quantity, 10)
			}

			return utils.PrepareErrorResponse(c, msg, consts.CErrorWalletNetworkIsNotReady, http.StatusServiceUnavailable)
		}

		return next(c)
	}
}

// ################################################################################
// GetNetworkInfo - function returns network info
// GetNetworkInfo godoc
// @Summary Get network info
// @Description Get network info
// @Tags Network
// @Produce  json
// @Success 200 {object} models.NetworkInfoResponse
// @Router /api/v1/network/info [get]
func (api *NetworkAPI) GetNetworkInfo(c echo.Context) error {
	return utils.PrepareSuccessResponse(c, api.walletNetworkInfo)
}
