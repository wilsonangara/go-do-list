// Package postgres provide storage level implementations
// to manage our database.
package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
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
		if err := runMigrations(db, d.MigrationDir); err != nil {
			return nil, err
		}
	}
	return &Storage{db: db}, nil
}

// NewTestStorage creates a new database for testing purposes.
func NewTestStorage(dbString, migrationDir string) (*sql.DB, func()) {
	db, teardown := newTestDB(dbString)

	db.SetMaxOpenConns(50)
	db.SetConnMaxLifetime(time.Minute)

	// run migrations when migration directory path is provided.
	if migrationDir != "" {
		if err := runMigrations(db, migrationDir); err != nil {
			// setting this to panic as we cannot run the tests when
			// migration fails.
			panic(err)
		}
	}

	return db, teardown
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

// runMigrations runs all the migration files in the path of migration directory
// specified.
func runMigrations(db *sql.DB, migrationDir string) error {
	if err := goose.SetDialect(driver); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}
	if err := goose.Up(db, migrationDir); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}
	return nil
}

// newTestDB creates a new isolated database for unit testing, on teardown,
// it will perform a cleanup by dropping the database entirely.
func newTestDB(ddlConnStr string) (*sql.DB, func()) {
	testDBName := uuid.New().String()

	// connect to database with the default connection string to create
	// a new database.
	ddlDB, err := sql.Open(driver, ddlConnStr)
	if err != nil {
		panic(fmt.Errorf("failed to open ddl database: %w", err))
	}
	ddlDB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, testDBName))
	if err := ddlDB.Close(); err != nil {
		panic(fmt.Errorf("failed to close ddl database: %w", err))
	}

	// connect to the newly created database
	connStr := replaceDBName(ddlConnStr, testDBName)
	db, err := sql.Open(driver, connStr)
	if err != nil {
		panic(fmt.Errorf("failed to open database: %w", err))
	}

	teardownFn := func() {
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close database: %v", err)
		}

		// connect to ddl database to drop test database.
		ddlDB, err := sql.Open(driver, ddlConnStr)
		if err != nil {
			log.Fatalf("failed to open ddl database: %v", err)
		}
		if _, err := ddlDB.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, testDBName)); err != nil {
			log.Fatalf("failed to drop test database: %v", err)
		}
		if err := ddlDB.Close(); err != nil {
			log.Fatalf("failed to close ddl database connection: %v", err)
		}
	}

	return db, teardownFn
}

// replaceDBName replaces the current connection string's dbname with the a new
// dbname.
func replaceDBName(dbString, dbName string) string {
	splitDBString := strings.Split(dbString, " ")

	newDBString := []string{}
	for _, v := range splitDBString {
		if strings.HasPrefix(v, "dbname") {
			newDBString = append(newDBString, fmt.Sprintf("dbname=%s", dbName))
			continue
		}
		newDBString = append(newDBString, v)
	}

	return strings.Join(newDBString, " ")
}
