package repository

import (
	"context"
	"log"

	walletPB "github.com/intellisoftalpin/proto/proto-gen/wallet"
	"github.com/intellisoftalpin/transactions-proxy-backend/models"
)

type PoolsRepo struct {
	config       *models.Config
	WalletClient walletPB.WalletClient
}

func NewPoolsRepo(config *models.Config, walletClient walletPB.WalletClient) *PoolsRepo {
	return &PoolsRepo{
		config:       config,
		WalletClient: walletClient,
	}
}

func (p *PoolsRepo) GetAllPools() (models.Pools, error) {
	pools := models.Pools{}

	for _, poolID := range p.config.Pools {
		pools.Pools = append(pools.Pools, models.Pool{PoolID: poolID})
	}

	return pools, nil
}

func (p *PoolsRepo) DelegateToPool(delegateToPool models.DelegateToPoolRequest) (string, error) {
	ctx := context.Background()

	log.Println("DelegateToPool request: ", delegateToPool)

	resp, err := p.WalletClient.SubmitTransaction(ctx, &walletPB.SubmitTransactionRequest{Tx: delegateToPool.CBOR})
	if err != nil {
		return "", err
	}

	log.Println("DelegateToPool response: ", resp)

	return "", nil
}
