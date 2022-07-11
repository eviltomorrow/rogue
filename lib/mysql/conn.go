package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//
var (
	DSN     string
	MinOpen int = 5
	MaxOpen int = 10
	DB      *sql.DB
)

var DefaultConnectTimeout = 10 * time.Second

// Build build mysql
func Build() error {
	pool, err := buildMySQL(DSN)
	if err != nil {
		return err
	}
	DB = pool
	return nil
}

// Close close mysql
func Close() error {
	if DB == nil {
		return nil
	}

	return DB.Close()
}

func buildMySQL(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DSN no set")
	}
	pool, err := sql.Open("mysql", DSN)
	if err != nil {
		return nil, err
	}
	pool.SetConnMaxLifetime(time.Minute * 3)
	pool.SetMaxOpenConns(MaxOpen)
	pool.SetMaxIdleConns(MinOpen)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultConnectTimeout)
	defer cancel()

	if err = pool.PingContext(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}

// Exec exec mysql
type Exec interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
