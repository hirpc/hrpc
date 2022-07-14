package mysql

type Options struct {
	Address  string `json:"address"`
	DBName   string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`

	// customized
	customized bool

	// MaxOpenConns sets the maximum number of open connections to the database
	MaxOpenConns int `json:"max_open_conns"`
	// MaxIdleConns sets the maximum number of connections in the idle connection pool.
	MaxIdleConns int `json:"max_idle_conns"`
}

type Option func(o *Options)

// WithCustomized will use your own configurations.
// To be reminder that you should make sure the values of Address, DB, Auth, Port have been assigned correctly.
func WithCustomized() Option {
	return func(o *Options) {
		o.customized = true
	}
}

func WithAddress(s string) Option {
	return func(o *Options) {
		o.Address = s
	}
}

func WithDB(name string) Option {
	return func(o *Options) {
		o.DBName = name
	}
}

func WithAuth(username, password string) Option {
	return func(o *Options) {
		o.Username = username
		o.Password = password
	}
}

func WithPort(port int) Option {
	return func(o *Options) {
		o.Port = port
	}
}

func WithMaxOpenConns(i int) Option {
	return func(o *Options) {
		o.MaxOpenConns = i
	}
}

func WithMaxIdleConns(i int) Option {
	return func(o *Options) {
		o.MaxIdleConns = i
	}
}
