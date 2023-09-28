package models

type Pools struct {
	Pools []Pool `json:"pools"`
}

type Pool struct {
	ID         int    `json:"id"`
	Ticker     string `json:"ticker"`
	Name       string `json:"name"`
	PoolID     string `json:"poolId"`     // BECH 32 Pool Id
	Saturation string `json:"saturation"` // Reward Stake in Lovelaces
	Pledge     string `json:"pledge"`     // Committed Pledge in Lovelaces
	Fee        string `json:"fee"`
	ROSe12     string `json:"rose12"`
}
