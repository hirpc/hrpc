package mysql

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
