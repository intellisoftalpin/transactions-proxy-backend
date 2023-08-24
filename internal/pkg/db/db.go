package db

import (
	"database/sql"
	"fmt"

	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/intellisoftalpin/transactions-proxy-backend/internal/pkg/db/migrations"
	"github.com/intellisoftalpin/transactions-proxy-backend/models"
)

// SetupDB - function to set up the PostgreSQL database
func SetupDB(dbConfig models.DBConfig) (db *sql.DB, err error) {
	// Open database connection
	pgsqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Database)
	db, err = sql.Open("postgres", pgsqlInfo)
	if err != nil {
		return nil, err
	}

	goose.SetBaseFS(migrations.Migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("failed to set up migrator: %w", err)
	}

	var opts []goose.OptionsFunc
	opts = append(opts, goose.WithAllowMissing())

	log.Println("running database migrations")

	if err := goose.Up(db, "sql", opts...); err != nil {
		log.Fatal("failed to run migrations: %w", err)
	}

	log.Println("database is ready")

	return db, nil
}
