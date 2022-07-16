package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	driver "github.com/go-sql-driver/mysql"
	"github.com/hirpc/hrpc/database"
	"github.com/hirpc/hrpc/uerror"
)

type mySQL struct {
	conn    *sql.DB
	options Options
}

// TxFunc can be used for transaction operation
// If error returned, tx.Rollback() will be called automatically
// If nil returned, tx.Commit() will be called automatically, also.
type TxFunc func(tx *sql.Tx) error

// NextFunc is designed for scanning all rows queryed from the database
// If error returned, it will cancel the loop in advance, and it will return the error.
// If ErrBreak returned, it will also cancel the loop in advance, but it will nil.
// If nil returned, it represents everything is OK.
type NextFunc func(*sql.Rows) error

// Proxy is a abstract layer for operating the MySQL database
type Proxy interface {
	// Transaction will start a transaction for the database
	// Ex.
	//		p.Transaction(ctx, func(tx *sql.Tx) error {
	//			tx.Exec(xxx)
	//			return nil // it will commit the transaction automatically
	//		})
	Transaction(ctx context.Context, fn TxFunc) error
	Query(ctx context.Context, query string, next NextFunc, args ...interface{}) error
	// QueryRow will query a row from the database
	// Ex.
	// 		var v1 string
	//		if err := p.QueryRow(ctx, `SELECT xx FROM xx`, []interface{}{&v1}, args); err != nil {
	//			// error
	//		}
	QueryRow(ctx context.Context, query string, dest []interface{}, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func (m mySQL) Transaction(ctx context.Context, fn TxFunc) error {
	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		return tx.Rollback()
	}
	tx.Commit()
	return nil
}

func (m mySQL) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.conn.ExecContext(ctx, query, args...)
}

var ErrBreak = uerror.New(300, "break")

func (m mySQL) Query(ctx context.Context, query string, next NextFunc, args ...interface{}) error {
	stmt, err := m.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := next(rows); err != nil {
			if err == ErrBreak {
				break
			}
			return err
		}
	}
	return nil
}

func (m mySQL) QueryRow(ctx context.Context, query string, dest []interface{}, args ...interface{}) error {
	stmt, err := m.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(dest...); err != nil {
		return err
	}
	return nil
}

var mm *mySQL

// Client returns the handler to operate mysql if success
func Client() Proxy {
	return mm
}

func (m *mySQL) Load(src []byte) error {
	// If the value of customized is true (enabled),
	// which means DOES NOT use the configurations from the configuration center.
	if m.options.customized {
		return nil
	}
	if err := json.Unmarshal(src, &m.options); err != nil {
		return err
	}
	return nil
}

func (m mySQL) dataSource() string {
	cfg := driver.Config{
		Addr:                    fmt.Sprintf("%s:%d", m.options.Address, m.options.Port),
		User:                    m.options.Username,
		Passwd:                  m.options.Password,
		Net:                     "tcp",
		DBName:                  m.options.DBName,
		AllowNativePasswords:    true,
		AllowCleartextPasswords: true,
	}
	return cfg.FormatDSN()
}

func (m *mySQL) Connect() error {
	m.Destory()

	db, err := sql.Open("mysql", m.dataSource())
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	db.SetMaxOpenConns(m.options.MaxOpenConns)
	db.SetMaxIdleConns(m.options.MaxIdleConns)
	m.conn = db
	return nil
}

// Valid returns a bool valud to determine whether the connection is ready to use
func Valid() bool {
	if mm == nil {
		return false
	}
	if mm.conn == nil {
		return false
	}
	if err := mm.conn.Ping(); err != nil {
		return false
	}
	return true
}

func (m mySQL) Name() string {
	return "mysql"
}

func (m *mySQL) Destory() {
	if m.conn != nil {
		m.conn.Close()
	}
}

func New(opts ...Option) *mySQL {
	var options = Options{
		Port:         3306,
		MaxOpenConns: 3,
		MaxIdleConns: 1,
		customized:   false,
	}
	for _, o := range opts {
		o(&options)
	}
	if mm != nil {
		mm.Destory()
	}
	mm = &mySQL{
		options: options,
	}
	return mm
}

var _ database.Database = (*mySQL)(nil)
