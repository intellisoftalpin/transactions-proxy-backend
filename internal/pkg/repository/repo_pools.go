package repository

import "github.com/intellisoftalpin/transactions-proxy-backend/models"

type PoolsRepo struct{}

func NewPoolsRepo() *PoolsRepo {
	return &PoolsRepo{}
}

func (p *PoolsRepo) GetAllPools() ([]models.Pool, error) {
	pools := []models.Pool{
		{
			ID:         1,
			Ticker:     "JNGL",
			Name:       "Jungle",
			PoolID:     "pool1lk6cxjaqd66t4t74q4gd9hymxapd93fvchhxt0uxwwprk9m8v6c",
			Saturation: "39403197265",
			Pledge:     "110000000",
			Fee:        "1%",
			ROSe12:     "2.25%",
		},
	}

	return pools, nil
}
