package storage

import (
	"fmt"
	"strings"
	"testing"

	"github.com/adepte-myao/avito_internship/internal/config"
	"github.com/sirupsen/logrus"
)

func TestStore(t *testing.T, databaseURL string) (*Storage, func()) {
	t.Helper()

	logger := logrus.New()
	logger.Level = logrus.ErrorLevel
	storeConfig := &config.StoreConfig{
		DatabaseURL: databaseURL,
	}
	storage := NewStorage(storeConfig, logger)
	if err := storage.Open(); err != nil {
		t.Fatal(err)
	}

	tables := []string{"accounts", "reserves_history"}
	if _, err := storage.db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))); err != nil {
		t.Fatal(err)
	}

	storage.db.Exec(`INSERT INTO services (id, name)
    					SELECT series.series, concat('service_', series.series)
       						FROM generate_series(1, 50) as series`)

	return storage, func() {
		storage.Close()
	}
}
