package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/adepte-myao/avito_internship/internal/config"
	"github.com/sirupsen/logrus"
)

func TestStore(t *testing.T, databaseURL string) (*sql.DB, func()) {
	t.Helper()

	logger := logrus.New()
	logger.Level = logrus.ErrorLevel
	storageConfig := &config.StoreConfig{
		DatabaseURL: databaseURL,
	}
	db, err := NewPostgresDb(storageConfig, logger)
	if err != nil {
		t.Fatal(err)
	}

	tables := []string{"accounts", "reserves_history"}
	if _, err := db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))); err != nil {
		t.Fatal(err)
	}

	db.Exec(`INSERT INTO services (id, name)
				SELECT series.series, concat('service_', series.series)
					FROM generate_series(1, 50) as series`)

	return db, func() {
		ClosePostgresDb(db)
	}
}
