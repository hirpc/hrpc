package redis

import (
	"encoding/json"
	"fmt"

	redisv8 "github.com/go-redis/redis/v8"
	"github.com/hirpc/hrpc/database/category"
)

type Option struct {
	Address  string `json:"address"`
	DB       int    `json:"db"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Network  string `json:"network"`
	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int `json:"max_retries"`
}

type Redis struct {
	conn   *redisv8.Client
	option Option
}

var (
	r *Redis
)

// Get returns the handler to operate redis if success
func Get() *redisv8.Client {
	return r.conn
}

func (r *Redis) Load(src []byte) error {
	if err := json.Unmarshal(src, &r.option); err != nil {
		return err
	}
	return nil
}

func (r Redis) dataSource() *redisv8.Options {
	cfg := &redisv8.Options{
		Network:    "tcp",
		Addr:       fmt.Sprintf("%s:%d", r.option.Address, r.option.Port),
		Username:   r.option.Username,
		Password:   r.option.Password,
		DB:         r.option.DB,
		MaxRetries: r.option.MaxRetries,
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

func (r Redis) Category() category.Category {
	return category.Redis
}

func (r *Redis) Destory() {
	if r.conn != nil {
		r.conn.Close()
	}
}

func New() *Redis {
	if r != nil {
		r.Destory()
	}
	r = &Redis{
		option: Option{MaxRetries: 3},
	}
	return r
}
