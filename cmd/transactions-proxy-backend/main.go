package main

import (
	"database/sql"
	"fmt"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/intellisoftalpin/transactions-proxy-backend/docs"
	"github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/config"
	tlpsdb "github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/db"
	models "github.com/intellisoftalpin/transactions-proxy-backend/models"
	"github.com/intellisoftalpin/transactions-proxy-backend/pkg/api"
)

// @title Transactions proxy backend API
// @version 0.1.1
// @description This is a service proxy app between frontend React.js app with Wallet Extention and Cardano Node backend app.
func main() {
	// log.SetFormatter(&log.TextFormatter{
	// 	CallerPrettyfier: func(f *runtime.Frame) (string, string) {
	// 		return "", fmt.Sprintf("%s:%d", f.File, f.Line)
	// 	},
	// })
	// log.SetReportCaller(true)

	loadedConfig, err := config.LoadConfig()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Setup db connection
	db, err := tlpsdb.SetupDB(loadedConfig.DB)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Close db connection on finish
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Println(err.Error())
		} else {
			fmt.Println("Database closed successfully.")
		}
	}(db)

	// Setup sessions
	var sessions = &models.Sessions{SessionsMap: make(map[string]*models.Session)}

	conn, err := setupClientConn(loadedConfig)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Setup server
	e := setupServer(db, sessions, loadedConfig, conn)

	// Start server
	e.Logger.Fatal(e.Start(":" + loadedConfig.ServerPort))
}

// setupServer - internal function to setup sever and router
func setupServer(db *sql.DB, sessions *models.Sessions, loadedConfig *models.Config, conn *grpc.ClientConn) (e *echo.Echo) {
	e = echo.New()

	e.Use(middleware.Logger())
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}, error=${error}\n",
	// }))
	e.Use(middleware.Recover())

	e.Use(middleware.CORS())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	apiHandlers := api.NewAPI(db, loadedConfig, sessions, conn)

	// Startup emptifying function
	go apiHandlers.UsersAPI.UsersSessionsEmptify()

	// Network
	// Get network info

	g := e.Group("/api/v1")
	g.GET("/network/info", apiHandlers.NetworkAPI.GetNetworkInfo)

	// Users
	// Login user with data from JSON
	g.POST("/user/login", apiHandlers.UsersAPI.LoginUser)

	// Transactions
	transactions := g.Group("/transactions")
	transactions.Use(apiHandlers.NetworkAPI.MiddlewareNetworkReady)

	transactions.Use(apiHandlers.TransactionsAPI.Middleware)

	// Get all user`s transactions
	transactions.GET("", apiHandlers.TransactionsAPI.GetAllTransactions)

	// Get active user`s transactions
	transactions.GET("/active", apiHandlers.TransactionsAPI.CheckActiveTransactions)

	// Get user`s single transaction
	transactions.GET("/:transaction_id", apiHandlers.TransactionsAPI.GetSingleTransaction)
	// Get user`s single transaction status
	transactions.GET("/:transaction_id/status", apiHandlers.TransactionsAPI.GetSingleTransactionStatus)

	// Create user`s transaction data from JSON
	transactions.POST("", apiHandlers.TransactionsAPI.CreateTransaction)

	// Update user`s single transaction
	transactions.PUT("/:transaction_id", apiHandlers.TransactionsAPI.UpdateSingleTransaction)
	// Change user`s single transaction status with data from JSON
	transactions.PUT("/:transaction_id/status", apiHandlers.TransactionsAPI.ChangeSingleTransactionStatus)

	// Delete user`s single transaction
	transactions.DELETE("/:transaction_id", apiHandlers.TransactionsAPI.DeleteSingleTransaction)

	// Tokens
	tokens := g.Group("/tokens")
	tokens.Use(apiHandlers.NetworkAPI.MiddlewareNetworkReady)

	// Get all tokens
	tokens.GET("", apiHandlers.TokensAPI.GetAllTokens)
	// Get single token
	tokens.GET("/:token_id", apiHandlers.TokensAPI.GetSingleToken)
	// Get single token price
	tokens.GET("/:token_id/price", apiHandlers.TokensAPI.GetSingleTokenPrice)

	// Pools
	pools := g.Group("/pools")

	// Get all pools
	pools.GET("", apiHandlers.PoolsAPI.GetAllPools)

	// Delegate to pool
	pools.POST("/delegate", apiHandlers.PoolsAPI.DelegateToPool, apiHandlers.NetworkAPI.MiddlewareNetworkReady)

	// Proxy
	proxy := g.Group("/proxy")

	// Submit External Transaction
	proxy.POST("/transactions", apiHandlers.ProxyAPI.SubmitExternalTransaction, apiHandlers.NetworkAPI.MiddlewareNetworkReady)

	return e
}

func setupClientConn(config *models.Config) (*grpc.ClientConn, error) {
	// certificate, err := tls.LoadX509KeyPair(config.CertPath+"/client-cert.pem", config.CertPath+"/client-key.pem")
	// if err != nil {
	// 	return nil, fmt.Errorf("could not load client key pair: %s", err)
	// }

	// certPool := x509.NewCertPool()
	// ca, err := os.ReadFile(config.CertPath + "/ca-cert.pem")
	// if err != nil {
	// 	return nil, fmt.Errorf("could not read ca certificate: %s", err)
	// }

	// if ok := certPool.AppendCertsFromPEM(ca); !ok {
	// 	return nil, fmt.Errorf("failed to append ca certs")
	// }

	// _ = certificate

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		// grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		// 	ServerName:   config.HostName,
		// 	Certificates: []tls.Certificate{certificate},
		// 	RootCAs:      certPool,
		// })),
	}

	conn, err := grpc.Dial(config.CNodeAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}

	return conn, nil
}
