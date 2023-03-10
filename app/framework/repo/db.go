package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"notification-service/domain"
	"time"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

func NewPostgesql(c *DBConfig) (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)

	db, err = initTCPSQL(c)
	zap.S().Info("Configured db using init TCP SQL")
	if err != nil {
		return nil, fmt.Errorf("initTCPSQL: unable to connect: %v", err)
	}

	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetConnMaxIdleTime(time.Duration(c.MaxIdleConns))

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		zap.S().Info("retrying db ping: ", err)
		time.Sleep(time.Second)
	}

	return db, err

}

func initTCPSQL(c *DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)

	return sql.Open("postgres", dsn)
}

func NewDBConfig() DBConfig {
	return DBConfig{}
}

type DBConfig struct {
	Type         string `env:"DATABASE_TYPE"`
	Host         string `env:"POSTGRESQL_HOST"`
	ConnName     string `env:"INSTANCE_CONNECTION_NAME"`
	User         string `env:"POSTGRESQL_USER"`
	Password     string `env:"POSTGRESQL_PASS"`
	DBName       string `env:"POSTGRESQL_DBNAME"`
	Port         int    `env:"POSTGRESQL_PORT"`
	MaxOpenConns int    `env:"POSTGRESQL_MAX_OPEN_CONNS"`
	MaxIdleConns int    `env:"POSTGRESQL_MAX_IDLE_CONNS"`
}

// contextExecutor can perform SQL queries with context
type contextExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// translate PSQL errors
func PSQLError(e error) (err error) {
	var pqErr *pq.Error
	if e != nil {
		if errors.As(e, &pqErr) && pqErr.Code[:2] == "08" { // connection error - should be retryable
			err = fmt.Errorf("%w: %v", domain.ErrInternal, e)
		} else {
			err = e
		}
	}
	return
}

// SQLParams ...
type SQLParams struct {
	Params []interface{}
}

// NewSQLParams ... useful for converter []T to []interface
func NewSQLParams(opt ParamOption) *SQLParams {
	p := &SQLParams{}
	opt(p)
	return p
}

// ParamOption ...
type ParamOption func(s *SQLParams)

// StringParam ...
func StringParam(strs []string) ParamOption {
	return func(s *SQLParams) {
		p := make([]interface{}, 0, len(strs))
		for _, v := range strs {
			p = append(p, v)
		}
		s.Params = p
	}
}
