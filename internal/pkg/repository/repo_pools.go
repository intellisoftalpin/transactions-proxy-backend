package repository

import "github.com/intellisoftalpin/transactions-proxy-backend/models"

type PoolsRepo struct {
	config *models.Config
}

func NewPoolsRepo(config *models.Config) *PoolsRepo {
	return &PoolsRepo{
		config: config,
	}
}

func (p *PoolsRepo) GetAllPools() (models.Pools, error) {
	pools := models.Pools{}

	for _, poolID := range p.config.Pools {
		pools.Pools = append(pools.Pools, models.Pool{PoolID: poolID})
	}

	return pools, nil
}
