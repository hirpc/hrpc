package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"

	mySQL "github.com/go-sql-driver/mysql"
	"github.com/hirpc/hrpc/database/category"
)

type Option struct {
	Address  string `json:"address"`
	DBName   string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`

	// MaxOpenConns sets the maximum number of open connections to the database
	MaxOpenConns int `json:"max_open_conns"`
	// MaxIdleConns sets the maximum number of connections in the idle connection pool.
	MaxIdleConns int `json:"max_idle_conns"`
}

type MySQL struct {
	conn   *sql.DB
	option Option
}

var (
	mm *MySQL
)

// Get returns the handler to operate mysql if success
func Get() *sql.DB {
	return mm.conn
}

func (m *MySQL) Load(src []byte) error {
	if err := json.Unmarshal(src, &m.option); err != nil {
		return err
	}
	return nil
}

func (m MySQL) dataSource() string {
	cfg := mySQL.Config{
		Addr:                    fmt.Sprintf("%s:%d", m.option.Address, m.option.Port),
		User:                    m.option.Username,
		Passwd:                  m.option.Password,
		Net:                     "tcp",
		DBName:                  m.option.DBName,
		AllowNativePasswords:    true,
		AllowCleartextPasswords: true,
	}
	return cfg.FormatDSN()
}

func (m *MySQL) Connect() error {
	m.Destory()

	db, err := sql.Open("mysql", m.dataSource())
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	db.SetMaxOpenConns(m.option.MaxOpenConns)
	db.SetMaxIdleConns(m.option.MaxIdleConns)
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

func (m MySQL) Category() category.Category {
	return category.MySQL
}

func (m *MySQL) Destory() {
	if m.conn != nil {
		m.conn.Close()
	}
}

func New() *MySQL {
	if mm != nil {
		mm.Destory()
	}
	mm = &MySQL{
		option: Option{
			MaxOpenConns: 3,
			MaxIdleConns: 1,
		},
	}
	return mm
}
