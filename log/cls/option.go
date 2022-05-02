package cls

type Options struct {
	endpoint string
	topicID  string
	callback Handle
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

func WithCallBack(fn Handle) Option {
	return func(o *Options) {
		o.callback = fn
	}
}
