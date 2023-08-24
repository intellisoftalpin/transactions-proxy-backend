package repository

import (
	"context"

	"github.com/intellisoftalpin/transactions-proxy-backend/models"
	walletPB "gitlab.com/encryptoteam/createtoken/token-lib-proto/proto-gen/wallet"
)

type TokensRepo struct {
	WalletClient walletPB.WalletClient
}

func NewTokensRepo(walletClient walletPB.WalletClient) *TokensRepo {
	return &TokensRepo{
		WalletClient: walletClient,
	}
}

func (t *TokensRepo) GetAllTokens() ([]models.Token, error) {
	ctx := context.Background()

	tokens, err := t.WalletClient.GetAllTokens(ctx, &walletPB.Empty{})
	if err != nil {
		return nil, err
	}

	return convertTokens(tokens.Tokens), nil
}

func (t *TokensRepo) GetSingleToken(tokenID string) (token models.Token, err error) {
	ctx := context.Background()

	resp, err := t.WalletClient.GetToken(ctx, &walletPB.TokenID{TokenId: tokenID})
	if err != nil {
		return token, err
	}

	return convertToken(resp.Token), err
}

func (t *TokensRepo) GetSingleTokenPrice(tokenID string) (tokenPrice models.TokenPrice, err error) {
	ctx := context.Background()

	resp, err := t.WalletClient.GetTokenPrice(ctx, &walletPB.TokenID{TokenId: tokenID})
	if err != nil {
		return tokenPrice, err
	}

	return convertTokenPrice(resp.Price), err
}

// ------------------------------------------------------------------------

func convertTokens(tokens []*walletPB.Token) []models.Token {
	var tokensResponse []models.Token

	for _, token := range tokens {
		tokensResponse = append(tokensResponse, convertToken(token))
	}

	return tokensResponse
}

func convertToken(token *walletPB.Token) models.Token {
	return models.Token{
		AssetName: token.AssetName,
		PolicyId:  token.PolicyId,
		AssetId:   token.AssetId,
		Ticker:    token.Ticker,
		Logo:      token.Logo,
		// Decimals:  token.Decimals,

		Address: token.Address,
		Price:   models.TokenPrice{Price: token.Price.Price},

		AssetUnit:     token.AssetUnit,
		AssetQuantity: token.AssetQuantity,
		Decimals:      token.AssetDecimals,
		// AssetDecimals: token.AssetDecimals,
		Fee:           token.Fee,
		Deposit:       token.Deposit,
		ProcessingFee: token.ProcessingFee,
		TotalQuantity: token.TotalQuantity,
		RewardAddress: token.RewardAddress,
	}
}

func convertTokenPrice(token *walletPB.Price) models.TokenPrice {
	return models.TokenPrice{
		Price: token.Price,
	}
}
