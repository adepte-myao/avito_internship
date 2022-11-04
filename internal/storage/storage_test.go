package storage_test

import (
	"os"
	"testing"
)

var databaseURL string

func TestMain(m *testing.M) {
	databaseURL = "postgres://balancer:superpass@localhost:5434/tests?sslmode=disable"

	os.Exit(m.Run())
}
