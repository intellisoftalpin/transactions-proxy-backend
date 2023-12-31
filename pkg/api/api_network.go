package api

import (
	"context"
	"log"
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
	cnodeNetworkInfo   models.NetworkInfoResponse

	walletsState []models.WalletState
}

func NewNetworkAPI(walletClient walletPB.WalletClient) *NetworkAPI {
	n := &NetworkAPI{
		WalletClient: walletClient,
	}

	go func(n *NetworkAPI) {
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

				n.cnodeNetworkInfo = models.ToNetworkInfoResponse(nInfo)
			}

			walletsState, err := n.WalletClient.GetWalletsState(ctx, &walletPB.Empty{})
			if err != nil {
				log.Println("GetWalletsState error: ", err.Error()) // 2023/10/10 17:49:20 GetWalletsState error:  rpc error: code = Unimplemented desc = unknown method GetWalletsState for service wallet.Wallet
				n.walletsState = []models.WalletState{}
				n.walletNetworkReady = false
			} else {
				n.walletsState = models.ToWalletsState(walletsState)

				for _, walletState := range n.walletsState {
					if walletState.Status != consts.CSyncProgressStatusReady {
						n.walletNetworkReady = false
						break
					}
				}
			}

		}
	}(n)

	return n
}

func (api *NetworkAPI) MiddlewareNetworkReady(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if !api.walletNetworkReady {
			msg := "Wallet network is not ready. Network Status: " + api.cnodeNetworkInfo.SyncProgress.Status

			if api.cnodeNetworkInfo.SyncProgress.Status == "syncing" {
				msg += ". Progress: " + strconv.FormatFloat(float64(api.cnodeNetworkInfo.SyncProgress.Progress.Quantity), 'f', 2, 32)
			}

			// не работает!!!!!!!!!!!!!!!
			for _, walletState := range api.walletsState {
				if walletState.Status != consts.CSyncProgressStatusReady {
					msg += ". Wallet Status: " + walletState.Status
					if walletState.Status == "syncing" {
						msg += ". Progress: " + strconv.FormatFloat(float64(walletState.Progress.Quantity), 'f', 2, 32)
					}
				}
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
	resp := struct {
		NetworkInfoResponse models.NetworkInfoResponse `json:"network_info"`
		WalletsState        []models.WalletState       `json:"wallets_state"`
	}{
		NetworkInfoResponse: api.cnodeNetworkInfo,
		WalletsState:        api.walletsState,
	}
	return utils.PrepareSuccessResponse(c, resp)
}
