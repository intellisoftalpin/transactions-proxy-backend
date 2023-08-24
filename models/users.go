package models

import (
	"sync"
	"time"
)

// Struct to store user`s data for login (and to insert new user into database if that user does not exist)
type CardanoUser struct {
	UserHash    string `json:"userHash" validate:"required|min_len:2"`
	UserRuntime int    `json:"userRuntime"`
}

// Struct to store user`s data (userID) after selecting user from database
type User struct {
	UserID uint64
}

// Struct to store user`s single session information
type Session struct {
	AuthorizationKey   string
	UserID             uint64
	ExpirationDateTime time.Time
}

// Struct to store all sessions information
type Sessions struct {
	Mux         sync.Mutex
	SessionsMap map[string]*Session
}
