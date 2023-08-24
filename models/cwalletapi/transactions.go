package cwalletapi

type Transaction struct {
	ID     string   `json:"id"`
	Amount Quantity `json:"amount"`
	Fee    Quantity `json:"fee"`
	// DepositTaken    Quantity `json:"deposit_taken"`
	// DepositReturned Quantity `json:"deposit_returned"`
	InsertedAt   Tip `json:"inserted_at"`
	ExpiresAt    Tip `json:"expires_at"`
	PendingSince Tip `json:"pending_since"`
	// Depth           Quantity `json:"depth"`
	// Direction       string    `json:"direction"`
	Inputs  []Input   `json:"inputs"`
	Outputs []Payment `json:"outputs"`
	// Collaterals     []Collateral `json:"collateral"`
	// CollateralOutputs []Payment    `json:"collateral_outputs"`
	// Withdrawals       []Withdrawal `json:"withdrawals"`
	Status         string   `json:"status"`
	Metadata       Metadata `json:"metadata"`
	ScriptValidity string   `json:"script_validity"`
	// Certificates      []Certificate `json:"certificates"`
	// Mint              Mint          `json:"mint"`
	// Burn              Burn          `json:"burn"`
	// ValidityInterval  ValidityInterval `json:"validity_interval"`
	// ScriptIntegrity   []string         `json:"script_integrity"`
	// ExtraSignatures   []string         `json:"extra_signatures"`
}

type Quantity struct {
	Quantity uint64 `json:"quantity"`
	Unit     string `json:"unit"`
}

type Tip struct {
	AbsoluteSlotNumber uint64   `json:"absolute_slot_number"`
	SlotNumber         uint64   `json:"slot_number"`
	EpochNumber        uint64   `json:"epoch_number"`
	Time               string   `json:"time"`
	Height             Quantity `json:"height"`
}

type Input struct {
	Payment
	ID    string `json:"id"`
	Index uint64 `json:"index"`
}
type Payment struct {
	Address string   `json:"address"`
	Amount  Quantity `json:"amount"`
	Assets  []Asset  `json:"assets"`
	// DerivationPath []string `json:"derivation_path"`
}
type Asset struct {
	PolicyID  string `json:"policy_id"`
	AssetName string `json:"asset_name"`
	Quantity  uint64 `json:"quantity"`
}

type Collateral struct {
	Address string   `json:"address"`
	Amount  Quantity `json:"amount"`
	ID      string   `json:"id"`
	Index   uint64   `json:"index"`
}

type Withdrawal struct {
	StakeAddress string   `json:"stake_address"`
	Amount       Quantity `json:"amount"`
}

type Metadata map[string]MetadataValue

type MetadataValue struct {
	String string          `json:"string,omitempty"`
	Int    uint64          `json:"int,omitempty"`
	Bytes  string          `json:"bytes,omitempty"`
	List   []MetadataValue `json:"list,omitempty"`
	Map    []MetadataMap   `json:"map,omitempty"`
}

type MetadataMap struct {
	K MetadataValue `json:"k"`
	V MetadataValue `json:"v"`
}

type Certificate struct {
	CertificateType   string   `json:"certificate_type"`
	Pool              string   `json:"pool"`
	RewardAccountPath []string `json:"reward_account_path"`
}

type Mint struct {
	Tokens               []Token `json:"tokens"`
	WalletPolicyKeyHash  string  `json:"wallet_policy_key_hash"`
	WalletPolicyKeyIndex string  `json:"wallet_policy_key_index"`
}

type Token struct {
	PolicyID     string       `json:"policy_id"`
	PloicyScript PloicyScript `json:"ploicy_script"`
	Assets       []TokenAsset `json:"assets"`
}

type PloicyScript struct {
	ScriptType string    `json:"script_type"`
	Script     string    `json:"script"`
	Reference  Reference `json:"reference"`
}

type Reference struct {
	ID    string `json:"id"`
	Index uint64 `json:"index"`
}

type TokenAsset struct {
	AssetName   string `json:"asset_name"`
	Quantity    uint64 `json:"quantity"`
	Fingerprint string `json:"fingerprint"`
}

type Burn struct {
	Mint
}

type ValidityInterval struct {
	InvalidBefore    Quantity `json:"invalid_before"`
	InvalidHereafter Quantity `json:"invalid_hereafter"`
}
