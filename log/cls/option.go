package cls

type Options struct {
	endpoint string
}

type Option func(*Options)

func WithEndpoint(endpoint string) Option {
	return func(o *Options) {
		o.endpoint = endpoint
	}
}
