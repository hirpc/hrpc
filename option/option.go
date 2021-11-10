package option

import (
	"errors"

	"github.com/hirpc/hrpc/database"
	"github.com/hirpc/hrpc/database/category"
	"github.com/hirpc/hrpc/life"
)

type Environment string

func (e Environment) String() string {
	return string(e)
}

const (
	Development Environment = "development"
	Production  Environment = "production"
)

var (
	ErrMissingServerName = errors.New("missing server name")
	ErrInvalidPort       = errors.New("invalid port number")
)

// Option defines
type Options struct {
	// ID is the service ID
	ID           string
	ServerName   string
	ListenPort   int
	ENV          Environment
	ConsulCenter Consul
	DBs          map[category.Category]database.Database
	HealthCheck  bool

	WhenExit    []life.Listener
	WhenRestart []life.Listener
}

func (o Options) Valid() error {
	if o.ServerName == "" {
		return ErrMissingServerName
	}
	if o.ListenPort == 0 {
		return ErrInvalidPort
	}
	return nil
}

type Option func(*Options)

// WithEnvironment sets the env
func WithEnvironment(env Environment) Option {
	return func(o *Options) {
		o.ENV = env
	}
}

func WithServerName(name string) Option {
	return func(o *Options) {
		o.ServerName = name
	}
}

func WithHealthCheck(v bool) Option {
	return func(o *Options) {
		o.HealthCheck = v
	}
}

func WithListenPort(port int) Option {
	return func(o *Options) {
		o.ListenPort = port
	}
}

func WithConsul(c Consul) Option {
	return func(o *Options) {
		o.ConsulCenter = c
	}
}

func WithDatabase(dbs ...database.Database) Option {
	return func(o *Options) {
		for _, db := range dbs {
			o.DBs[db.Category()] = db
		}
	}
}

func WithExitListeners(listeners ...life.Listener) Option {
	return func(o *Options) {
		o.WhenExit = append(o.WhenExit, listeners...)
	}
}

func WithRestartListeners(listeners ...life.Listener) Option {
	return func(o *Options) {
		o.WhenRestart = append(o.WhenRestart, listeners...)
	}
}
