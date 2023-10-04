package config

import (
	"os"
	"strings"

	"github.com/intellisoftalpin/transactions-proxy-backend/models"
)

func LoadConfig() (loadedConfig *models.Config, err error) {
	loadedConfig = &models.Config{
		ServerPort:   os.Getenv("SERVER_PORT"),
		CNodeAddress: os.Getenv("CNODE_ADDRESS"),
		HostName:     os.Getenv("HOST_NAME"),
		CertPath:     os.Getenv("PATH_TO_CERTS"),
		DB: models.DBConfig{
			Host:     os.Getenv("POSTGRES_DB_HOST"),
			Port:     os.Getenv("POSTGRES_DB_PORT"),
			User:     os.Getenv("POSTGRES_DB_USER"),
			Password: os.Getenv("POSTGRES_DB_PASS"),
			Database: os.Getenv("POSTGRES_DB_NAME"),
		},
	}

	loadedConfig.Pools = strings.Split(os.Getenv("POOLS"), ";")

	// log.Infoln("Loaded config:", loadedConfig)

	return loadedConfig, nil
}
