package configs

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
)

func init() {
	if token := os.Getenv("CONFIGS_TOKEN"); token != "" {
		o.Token = token
	}
	if addr := os.Getenv("CONFIGS_ADDR"); addr != "" {
		o.Address = addr
	}
	if dc := os.Getenv("CONFIGS_DATACENTER"); dc != "" {
		o.DataCenter = dc
	}
}

var (
	c *api.Client
	o = Option{
		DataCenter: "dc1",
		Prefix:     "development",
	}
)

var (
	// ErrKeyNotExist for invalid key provided
	ErrKeyNotExist = errors.New("the key does not exist")
	// ErrNotInitialized for uninitialized
	ErrNotInitialized = errors.New("the client does not be initialized, call With() first")
)

// Option an option set
type Option struct {
	Token      string
	DataCenter string
	Address    string
	Prefix     string
}

func Register(opts ...Option) error {
	for _, opt := range opts {
		if opt.Token != "" {
			o.Token = opt.Token
		}
		if opt.DataCenter != "" {
			o.DataCenter = opt.DataCenter
		}
		if opt.Address != "" {
			o.Address = opt.Address
		}
		if opt.Prefix != "" {
			o.Prefix = opt.Prefix
		}
	}
	return client()
}

// Get returns the value of the key from consul
func Get(key string) ([]byte, error) {
	if c == nil {
		return nil, ErrNotInitialized
	}
	data, _, err := c.KV().Get(
		fmt.Sprintf("%s/%s", o.Prefix, key),
		nil,
	)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrKeyNotExist
	}
	return data.Value, nil
}

// Set will push a config pair to consul
func Set(key, val string) error {
	if c == nil {
		return ErrNotInitialized
	}
	_, err := c.KV().Put(
		&api.KVPair{
			Key:   fmt.Sprintf("%s/%s", o.Prefix, key),
			Value: []byte(val),
		}, nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func Prefix() string {
	return o.Prefix
}

// Client returns the consul client
func Client() *api.Client {
	return c
}

func client() error {
	config := api.DefaultConfig()
	config.Address = o.Address
	config.Token = o.Token
	config.Datacenter = o.DataCenter
	v, err := api.NewClient(config)
	if err != nil {
		return err
	}
	c = v
	return nil
}
