package api

import (
	"net/http"

	"github.com/intellisoftalpin/transactions-proxy-backend/consts"
	"github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/repository"
	utils "github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/utils"
	"github.com/labstack/echo/v4"
	walletPB "gitlab.com/encryptoteam/createtoken/token-lib-proto/proto-gen/wallet"
)

type TokensAPI struct {
	tokensRepo   *repository.TokensRepo
	WalletClient walletPB.WalletClient
}

func NewTokensAPI(walletClient walletPB.WalletClient) *TokensAPI {
	return &TokensAPI{
		tokensRepo:   repository.NewTokensRepo(walletClient),
		WalletClient: walletClient,
	}
}

// ################################################################################
// GetAllTokens - function to get all tokens
// GetAllTokens godoc
// @Summary Get All Tokens.
// @Description Get all tokens.
// @Tags tokens
// @Accept json
// @Produce json
// @Success 200 {object} []models.Token
// @Router /api/v1/tokens [GET]
func (api *TokensAPI) GetAllTokens(c echo.Context) error {
	tokens, err := api.tokensRepo.GetAllTokens()
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, tokens)
}

// ################################################################################
// GetSingleToken - function to get single token
// GetSingleToken godoc
// @Summary Get Single Token.
// @Description Get single token.
// @Tags tokens
// @Accept json
// @Produce json
// @Success 200 {object} models.Token
// @Router /api/v1/tokens/{token_id} [GET]
func (api *TokensAPI) GetSingleToken(c echo.Context) error {
	token, err := api.tokensRepo.GetSingleToken(c.Param("token_id"))
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, token)
}

// ################################################################################
// GetSingleTokenPrice - function to get single token price
// GetSingleTokenPrice godoc
// @Summary Get Single Token Price.
// @Description Get single token price.
// @Tags tokens
// @Accept json
// @Produce json
// @Success 200 {object} models.TokenPrice
// @Router /api/v1/tokens/{token_id}/price [GET]
func (api *TokensAPI) GetSingleTokenPrice(c echo.Context) error {
	token, err := api.tokensRepo.GetSingleTokenPrice(c.Param("token_id"))
	if err != nil {
		return utils.PrepareErrorResponse(c, err.Error(), consts.CErrorsInternalError, http.StatusInternalServerError)
	}

	return utils.PrepareSuccessResponse(c, token)
}
