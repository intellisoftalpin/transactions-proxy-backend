package models

type Pools struct {
	Pools []Pool `json:"pools"`
}

type Pool struct {
	ID         int    `json:"id,omitempty"`
	Ticker     string `json:"ticker,omitempty"`
	Name       string `json:"name,omitempty"`
	PoolID     string `json:"poolId"`               // BECH 32 Pool Id
	Saturation string `json:"saturation,omitempty"` // Reward Stake in Lovelaces
	Pledge     string `json:"pledge,omitempty"`     // Committed Pledge in Lovelaces
	Fee        string `json:"fee,omitempty"`
	ROSe12     string `json:"rose12,omitempty"`
}
