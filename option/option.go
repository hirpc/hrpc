package option

import (
	"errors"

	"github.com/hirpc/hrpc/database"
	"github.com/hirpc/hrpc/database/category"
	"github.com/hirpc/hrpc/life"
	"github.com/hirpc/hrpc/plugin"
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

	// StackSkip for logging that it can be used to debug stacks
	// default: 1
	StackSkip int
	Plugins   []plugin.Plugin

	// ServerCerts ...
	ServerCerts *certs

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

func WithStackSkip(i int) Option {
	return func(o *Options) {
		o.StackSkip = i
	}
}

func WithServerName(name string) Option {
	return func(o *Options) {
		o.ServerName = name
	}
}

func WithHealthCheck() Option {
	return func(o *Options) {
		o.HealthCheck = true
	}
}

func WithListenPort(port int) Option {
	return func(o *Options) {
		o.ListenPort = port
	}
}

func WithDatabases(dbs ...database.Database) Option {
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
