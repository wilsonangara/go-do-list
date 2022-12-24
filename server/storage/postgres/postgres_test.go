package postgres

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var dbString string

func TestMain(m *testing.M) {
	dbString = os.Getenv("DATABASE_CONNECTION")
	os.Exit(m.Run())
}

func TestNewStorage(t *testing.T) {
	t.Parallel()

	if dbString == "" {
		t.Skip("database connection not set, skipping.")
	}
	dbInfo := parseDBString(dbString)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		_, err := NewStorage(dbInfo)
		if err != nil {
			t.Fatalf("NewStorage(_) expected nil error, got = %v", err)
		}
	})
}

func TestNewTestStorage(t *testing.T) {
	t.Parallel()

	if dbString == "" {
		t.Skip("database connection not set, skipping.")
	}

	t.Run("Success_WithoutMigration", func(t *testing.T) {
		t.Parallel()

		_, teardown := NewTestStorage(dbString, "")
		t.Cleanup(teardown)
	})

	t.Run("Success_WithMigration", func(t *testing.T) {
		t.Parallel()

		var migrationDir = filepath.Join("..", "migrations")

		_, teardown := NewTestStorage(dbString, migrationDir)
		t.Cleanup(teardown)
	})
}

func parseDBString(dbString string) *DBInfo {
	splitDBString := strings.Split(dbString, " ")

	res := &DBInfo{
		SSLMode: "require",
	}

	for _, v := range splitDBString {
		splitValue := strings.Split(v, "=")
		value := splitValue[1]
		switch splitValue[0] {
		case "user":
			res.User = value
		case "password":
			res.Password = value
		case "host":
			res.Host = value
		case "port":
			res.Port = value
		case "dbname":
			res.Name = value
		case "connect_timeout":
			res.ConnectTimeout = value
		case "sslmode":
			res.SSLMode = value
		}
	}

	return res
}
