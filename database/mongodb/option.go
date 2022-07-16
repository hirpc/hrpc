package mongodb

type Options struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`

	// customized
	customized bool
}

type Option func(o *Options)

// WithCustomized will use your own configurations.
// To be reminder that you should make sure the values of Address, Auth, Port have been assigned correctly.
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

func WithAuth(username, password string) Option {
	return func(o *Options) {
		o.Username = username
		o.Password = password
	}
}
