package option

type Environment string

func (e Environment) String() string {
	return string(e)
}

const (
	Development Environment = "development"
	Production  Environment = "production"
)

// WithEnvironment sets the env
func WithEnvironment(env Environment) Option {
	return func(o *Options) {
		o.ENV = env
	}
}
