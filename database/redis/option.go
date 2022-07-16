package redis

type Options struct {
	Address  string `json:"address"`
	DB       int    `json:"db"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Network  string `json:"network"`
	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int `json:"max_retries"`

	// customized
	customized bool
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

func WithDB(i int) Option {
	return func(o *Options) {
		o.DB = i
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

func WithMaxRetries(i int) Option {
	return func(o *Options) {
		o.MaxRetries = i
	}
}
