package storage

import (
	"database/sql"

	"github.com/adepte-myao/avito_internship/internal/config"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func NewPostgresDb(config *config.StoreConfig, logger *logrus.Logger) (*sql.DB, error) {
	logger.Info("Connecting to database...")

	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		logger.Fatal("Connection to database was not successful. More: ", err.Error())
		return nil, err
	}
	logger.Info("Connected to database")

	if err := db.Ping(); err != nil {
		logger.Fatal("Ping to database was not successful. More: ", err.Error())
		return nil, err
	}
	logger.Info("Ping to database is successful")

	return db, nil
}

func ClosePostgresDb(db *sql.DB) {
	db.Close()
}
