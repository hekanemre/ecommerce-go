package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Repository defines the interface for database operations
type Repository interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Begin() (*sql.Tx, error)
	Close() error
}

// mysqlRepository implements Repository
type mysqlRepository struct {
	db *sql.DB
}

// NewMySQLRepository connects to MySQL and returns a repository
func NewMySQLRepository(user, password, host, port, dbname string) (Repository, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &mysqlRepository{db: db}, nil
}

// QueryRow executes a query that returns at most one row
func (r *mysqlRepository) QueryRow(query string, args ...interface{}) *sql.Row {
	return r.db.QueryRow(query, args...)
}

// Query executes a query that returns multiple rows
func (r *mysqlRepository) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return r.db.Query(query, args...)
}

// Exec executes a query without returning rows
func (r *mysqlRepository) Exec(query string, args ...interface{}) (sql.Result, error) {
	return r.db.Exec(query, args...)
}

// Begin starts a new transaction
func (r *mysqlRepository) Begin() (*sql.Tx, error) {
	return r.db.Begin()
}

// Close closes the database connection
func (r *mysqlRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return errors.New("db connection is nil")
}
