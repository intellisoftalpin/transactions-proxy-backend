package db

import (
	"database/sql"

	models "github.com/intellisoftalpin/transactions-proxy-backend/models"
)

type UsersDB struct {
	db *sql.DB
}

func NewUsersDB(db *sql.DB) (userDB *UsersDB) {
	return &UsersDB{db: db}
}

// Select query to get existed user`s ID
const querySelectUser = `
	SELECT id
	FROM users
	WHERE user_hash = $1 AND user_runtime = $2
`

// SelectUser - function to select user`s ID from database table if that user does exist
func (t *UsersDB) SelectUser(userHash string, userRuntime int) (user models.User, err error) {
	row := t.db.QueryRow(querySelectUser, userHash, userRuntime)

	if err = row.Scan(&user.UserID); err != nil {
		return user, err
	}

	return user, nil
}

// Insert new user query
const queryInsertUser = `
	INSERT INTO users (
		user_hash,
		user_runtime
	)
	VALUES (
		$1,
		$2
	)
`

// CreateUser - function to create new user by User Hash and User Runtime variable in database table
func (t *UsersDB) CreateUser(userHash string, userRuntime int) (err error) {
	query, err := t.db.Prepare(queryInsertUser)
	if err != nil {
		return err
	}

	if _, err = query.Exec(userHash, userRuntime); err != nil {
		return err
	}

	return nil
}
