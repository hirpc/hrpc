package option

type Consul struct {
	Address    string
	Token      string
	DataCenter string
	// Prefix default value will be the current environment's name
	Prefix string
}

func WithConsul(c Consul) Option {
	return func(o *Options) {
		o.ConsulCenter = c
	}
}
