package models

// ErrorResponse - error response structure
type ErrorResponse struct {
	Status    string `json:"status" example:"error"`
	MsgText   string `json:"msg" example:"Unknown error"`
	ErrorCode string `json:"errorCode" example:"unknownError"`
}

// Struct to store response about user`s session data with Authorization Key and Expiration Date Time if session is session not expired
type UserSessionResponse struct {
	Status             string `json:"status"`
	Message            string `json:"message"`
	AuthorizationKey   string `json:"sessionAuthorizationKey"`
	ExpirationDateTime string `json:"expirationDateTime"`
}

// --------------------------------------------------------------------------------

// Struct to store user`s single transaction response
type SingleTransactionResponse struct {
	Transaction
	TransactionData
}

// Struct to store transaction ID after inserting or updating user`s single transaction data
type SaveTransactionDataResponse struct {
	TransactionID uint64 `json:"transactionId"`
}

// Struct to store user`s single transaction status response
type SingleTransactionStatusResponse struct {
	TransactionStatus string `json:"transactionStatus"`
}

// Struct to store new transaction status after user`s single transaction status
type ChangeSingleTransactionStatusResponse struct {
	Status string `json:"status"`
}

// Struct to store response rezults of deliting user`s single transaction
type DeleteSingleTransactionResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
