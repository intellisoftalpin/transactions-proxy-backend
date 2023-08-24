package models

type Token struct {
	AssetName string `json:"assetName"`
	PolicyId  string `json:"policyId"`
	AssetId   string `json:"assetId"`
	Ticker    string `json:"ticker"`
	Logo      string `json:"logo"`
	Decimals  uint64 `json:"decimals"`

	Address string     `json:"address"`
	Price   TokenPrice `json:"tokenPrice"`

	AssetUnit     string `json:"assetUnit"`
	AssetQuantity uint64 `json:"assetQuantity"`
	Fee           uint64 `json:"fee"`
	Deposit       uint64 `json:"deposit"`
	ProcessingFee uint64 `json:"processingFee"`
	TotalQuantity uint64 `json:"totalQuantity"`
	RewardAddress string `json:"rewardAddress"`
}

type Tokens struct {
	Tokens []Token `json:"tokens"`
}

type TokenPrice struct {
	Price uint64 `json:"price"`
}
