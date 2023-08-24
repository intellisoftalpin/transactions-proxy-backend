package api

import (
	"database/sql"

	"google.golang.org/grpc"

	models "github.com/intellisoftalpin/transactions-proxy-backend/models"
	walletPB "gitlab.com/encryptoteam/createtoken/token-lib-proto/proto-gen/wallet"
)

type API struct {
	UsersAPI        *UsersAPI
	TransactionsAPI *TransactionsAPI
	TokensAPI       *TokensAPI

	Config *models.Config
}

func NewAPI(db *sql.DB, config *models.Config, sessions *models.Sessions, conn *grpc.ClientConn) *API {
	walletClient := walletPB.NewWalletClient(conn)

	return &API{
		UsersAPI:        NewUsersAPI(db, config, sessions),
		TransactionsAPI: NewTransactionsAPI(db, config, sessions, walletClient),
		TokensAPI:       NewTokensAPI(walletClient),

		Config: config,
	}
}
