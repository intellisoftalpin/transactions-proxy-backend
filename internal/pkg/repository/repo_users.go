package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/intellisoftalpin/transactions-proxy-backend/consts"
	tlpsdb "github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/db"
	models "github.com/intellisoftalpin/transactions-proxy-backend/models"
)

type UsersRepo struct {
	UsersDB       tlpsdb.UsersDBInterface
	UsersSessions *models.Sessions
}

func NewUsersRepo(db *sql.DB, config *models.Config, sessions *models.Sessions) (userRepo *UsersRepo) {
	return &UsersRepo{
		UsersDB:       tlpsdb.NewUsersDB(db),
		UsersSessions: sessions,
	}
}

// createNewSession - internal function to create new Authorization Key by user ID
func (t *UsersRepo) createNewSession(userID uint64, currTime time.Time) (authorizationKey string, expirTime string, err error) {
	authorizationKey = uuid.New().String()
	// authorizationKey = "1"
	expirationTime := currTime.Add(time.Hour * consts.SessionLifetimeInHours)
	t.UsersSessions.SessionsMap[authorizationKey] = &models.Session{
		AuthorizationKey:   authorizationKey,
		UserID:             userID,
		ExpirationDateTime: expirationTime,
	}

	return authorizationKey, expirationTime.String(), nil
}

// getSession - internal function to get Authorization Key by user ID if it is not expired
func (t *UsersRepo) getSession(userID uint64, currTime time.Time) (authorizationKey string, expirTime string, sessionMessage string, err error) {
	for k, v := range t.UsersSessions.SessionsMap {
		if userID == v.UserID {
			// User exists
			if currTime.Before(v.ExpirationDateTime) {
				// Session is not expired
				authorizationKey = k
				expirTime = v.ExpirationDateTime.String()
				sessionMessage = consts.MessageSessionNotExpiredOldKeySuccess
				return authorizationKey, expirTime, sessionMessage, nil
			}
			break
		}
	}

	// Session is expired or user exists but that user`s session does not exists at all. Create new session
	authorizationKey, expirTime, err = t.createNewSession(userID, currTime)
	if err != nil {
		sessionMessage = consts.MessageSessionNotExistsOrExpiredNewSessionError
		return authorizationKey, expirTime, sessionMessage, err
	}

	sessionMessage = consts.MessageSessionNotExistsOrExpiredNewSessionSuccess

	return authorizationKey, expirTime, sessionMessage, nil
}

// LoginUserSession - function to login (create new if needed) user by Hash and Runtime variable and return Authorization Key (Session ID)
func (t *UsersRepo) LoginUserSession(userHash string, userRuntime int) (userResponse *models.UserSessionResponse, err error) {
	currentTime := time.Now()

	userResponse = &models.UserSessionResponse{
		Status:  consts.CErrorStatus,
		Message: consts.MessageOtherError,
	}

	// Check if user exists
	user, err := t.UsersDB.SelectUser(userHash, userRuntime)
	if err != nil && err == sql.ErrNoRows {
		// Create new user if does not exists
		if err = t.UsersDB.CreateUser(userHash, userRuntime); err != nil {
			return userResponse, err
		}

		// Get new user`s ID
		user, err = t.UsersDB.SelectUser(userHash, userRuntime)
		if err != nil {
			return userResponse, err
		}

		//Create new session and get session ID
		userResponse.AuthorizationKey, userResponse.ExpirationDateTime, err = t.createNewSession(user.UserID, currentTime)
		if err != nil {
			userResponse.Message = consts.MessageNewUserAndNewSessionError
			return userResponse, err
		}

		userResponse.Status = consts.CSuccessStatus
		userResponse.Message = consts.MessageNewUserAndNewSessionSuccess
		return userResponse, nil
	} else if err != nil {
		return userResponse, err
	}

	// Get session ID or create new session if user exists
	var sessionMessage string
	userResponse.AuthorizationKey, userResponse.ExpirationDateTime, sessionMessage, err = t.getSession(user.UserID, currentTime)
	if err != nil {
		userResponse.Status = consts.CErrorStatus
		userResponse.Message = consts.MessageUserExists + sessionMessage
	}

	userResponse.Status = consts.CSuccessStatus
	userResponse.Message = consts.MessageUserExists + sessionMessage

	return userResponse, nil
}

// SessionsEmptify - function to emptify list of all sessions after a period of time
func (t *UsersRepo) SessionsEmptify() {
	for {
		time.Sleep(consts.SessionsEmptifyPeriodInHours * time.Hour)

		t.UsersSessions.Mux.Lock()

		fmt.Printf("Length before cleanup: %d\n", len(t.UsersSessions.SessionsMap))

		for k := range t.UsersSessions.SessionsMap {
			delete(t.UsersSessions.SessionsMap, k)
		}

		fmt.Printf("Length after cleanup: %d\n", len(t.UsersSessions.SessionsMap))

		t.UsersSessions.Mux.Unlock()
	}
}
