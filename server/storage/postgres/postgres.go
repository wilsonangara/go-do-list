// Package postgres provide storage level implementations
// to manage our database.
package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

const (
	driver = "postgres"
)

// Storage provides a wrapper on our database that provides
// required methods to interact with our data.
type Storage struct {
	db *sql.DB
}

func (s *Storage) DB() *sql.DB {
	return s.db
}

type DBInfo struct {
	User           string
	Host           string
	Name           string
	Port           string
	Password       string
	ConnectTimeout string
	SSLMode        string
	MigrationDir   string
}

// NewStorage creates a new database connection with the provided
// credentials.
func NewStorage(d *DBInfo) (*Storage, error) {
	dbString := d.dbStringBuilder()
	db, err := sql.Open(driver, dbString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// run migrations when migration directory path is provided.
	if d.MigrationDir != "" {
		if err := goose.SetDialect(driver); err != nil {
			return nil, fmt.Errorf("failed to set goose dialect: %w", err)
		}
		if err := goose.Up(db, d.MigrationDir); err != nil {
			return nil, fmt.Errorf("failed to run database migrations: %w", err)
		}
	}
	return &Storage{db: db}, nil
}

// dbStringBuilder builds a database connection string with the provided
func (d *DBInfo) dbStringBuilder() string {
	dbString := []string{}

	dbString = append(dbString, fmt.Sprintf("user=%s", d.User))
	dbString = append(dbString, fmt.Sprintf("host=%s", d.Host))
	dbString = append(dbString, fmt.Sprintf("port=%s", d.Port))
	dbString = append(dbString, fmt.Sprintf("dbname=%s", d.Name))
	dbString = append(dbString, fmt.Sprintf("connect_timeout=%s", d.ConnectTimeout))
	dbString = append(dbString, fmt.Sprintf("sslmode=%v", d.SSLMode))

	// check if the database has a password.
	if d.Password != "" {
		dbString = append(dbString, fmt.Sprintf("password=%s", d.Password))
	}

	return strings.Join(dbString, " ")
}
