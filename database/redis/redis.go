package redis

import (
	"encoding/json"
	"fmt"

	redisv8 "github.com/go-redis/redis/v8"
	"github.com/hirpc/hrpc/database"
)

type Redis struct {
	conn    *redisv8.Client
	options Options
}

var r *Redis

var (
	// ErrNil when key does not exist.
	ErrNil = redisv8.Nil
)

// Client returns the handler to operate redis if success
func Client() *redisv8.Client {
	return r.conn
}

func (r *Redis) Load(src []byte) error {
	// If the value of customized is true (enabled),
	// which means DOES NOT use the configurations from the configuration center.
	if r.options.customized {
		return nil
	}
	if err := json.Unmarshal(src, &r.options); err != nil {
		return err
	}
	return nil
}

func (r Redis) dataSource() *redisv8.Options {
	cfg := &redisv8.Options{
		Network:    "tcp",
		Addr:       fmt.Sprintf("%s:%d", r.options.Address, r.options.Port),
		Username:   r.options.Username,
		Password:   r.options.Password,
		DB:         r.options.DB,
		MaxRetries: r.options.MaxRetries,
	}
	return cfg
}

func (r *Redis) Connect() error {
	r.Destory()

	r.conn = redisv8.NewClient(r.dataSource())
	pong, err := r.conn.Ping(r.conn.Context()).Result()
	if err != nil || pong != "PONG" {
		return err
	}
	return nil
}

// Valid returns a bool valud to determine whether the connection is ready to use
func Valid() bool {
	if r == nil {
		return false
	}
	if r.conn == nil {
		return false
	}
	pong, err := r.conn.Ping(r.conn.Context()).Result()
	if err != nil || pong != "PONG" {
		return false
	}
	return true
}

func (r Redis) Name() string {
	return "redis"
}

func (r *Redis) Destory() {
	if r.conn != nil {
		r.conn.Close()
	}
}

func New(opts ...Option) *Redis {
	var options = Options{
		Port:       6379,
		DB:         0,
		MaxRetries: 3,
		customized: false,
	}
	for _, o := range opts {
		o(&options)
	}

	if r != nil {
		r.Destory()
	}
	r = &Redis{
		options: options,
	}
	return r
}

var _ database.Database = (*Redis)(nil)
