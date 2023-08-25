package api

import (
	"database/sql"

	"google.golang.org/grpc"

	walletPB "github.com/intellisoftalpin/proto/proto-gen/wallet"
	models "github.com/intellisoftalpin/transactions-proxy-backend/models"
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
