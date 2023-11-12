package repository

import (
	"context"
	"log"

	walletPB "github.com/intellisoftalpin/proto/proto-gen/wallet"
	"github.com/intellisoftalpin/transactions-proxy-backend/models"
)

type ProxyRepo struct {
	config       *models.Config
	WalletClient walletPB.WalletClient
}

func NewProxyRepo(config *models.Config, walletClient walletPB.WalletClient) *ProxyRepo {
	return &ProxyRepo{
		config:       config,
		WalletClient: walletClient,
	}
}

func (p *ProxyRepo) SubmitExternalTransaction(submitExternalTransaction models.SubmitExternalTransactionRequest) (txHash string, err error) {
	ctx := context.Background()

	log.Println("SubmitExternalTransaction request: ", submitExternalTransaction)

	resp, err := p.WalletClient.SubmitTransaction(ctx, &walletPB.SubmitTransactionRequest{Tx: submitExternalTransaction.CBOR})
	if err != nil {
		log.Println("SubmitExternalTransaction error: ", err)
		return "", err
	}

	log.Println("SubmitExternalTransaction response: ", resp)

	return resp.TxHash, nil
}
