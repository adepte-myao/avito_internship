package storage

import (
	"database/sql"
	"time"

	"github.com/adepte-myao/avito_internship/internal/config"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	config *config.StoreConfig
	db     *sql.DB
	logger *logrus.Logger
}

func NewStorage(config *config.StoreConfig, logger *logrus.Logger) *Storage {
	return &Storage{
		config: config,
		logger: logger,
	}
}

func (s *Storage) Open() error {
	s.logger.Info("Connecting to database: first attempt")

	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		// Sometimes database isn't up when we try to connect it.
		// Take some time, maybe db will be up
		s.logger.Info("Connecting to database: second attempt")

		time.Sleep(10 * time.Second)
		if db, err = sql.Open("postgres", s.config.DatabaseURL); err != nil {
			return err
		}
	}
	s.logger.Info("Connected to database")

	if err := db.Ping(); err != nil {
		s.logger.Error("Ping to database was not successful. Reason: ", err.Error())
		return err
	}
	s.logger.Info("Ping to database is successful")

	s.db = db

	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}
