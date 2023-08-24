package db

import models "github.com/intellisoftalpin/transactions-proxy-backend/models"

// Interface for users database model
type UsersDBInterface interface {
	SelectUser(userHash string, userRuntime int) (user models.User, err error)
	CreateUser(userHash string, userRuntime int) (err error)
}
