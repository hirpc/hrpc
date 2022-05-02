package cls

type Options struct {
	endpoint string
	topicID  string
}

type Option func(*Options)

func WithEndpoint(endpoint string) Option {
	return func(o *Options) {
		o.endpoint = endpoint
	}
}

func WithTopicID(id string) Option {
	return func(o *Options) {
		o.topicID = id
	}
}
