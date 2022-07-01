package configs

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
)

var (
	// ErrKeyNotExist for invalid key provided
	ErrKeyNotExist = errors.New("the key does not exist")
	// ErrNotInitialized for uninitialized
	ErrNotInitialized = errors.New("the client does not be initialized, call With() first")
)

type Configs interface {
	Prefix() string

	Get(key string) ([]byte, error)
	SGet(k string) ([]byte, error)

	Set(key, val string) error
}

type configs struct {
	Token      string
	DataCenter string
	Address    string
	prefix     string

	client *api.Client
}

var (
	cfg        configs
	serverName string
)

func New(addr, dc, prefix, token, sName string) *configs {
	serverName = sName
	cfg = configs{
		DataCenter: dc,
		prefix:     prefix,
		Address:    addr,
		Token:      token,
	}
	return &cfg
}

func (c *configs) DependsOn() []string {
	return nil
}

func (c *configs) Name() string {
	return "hrpc-configs"
}

func (c *configs) Loaded() bool {
	return c != nil && c.client != nil
}

func (c *configs) mergeENV() {
	if token := os.Getenv("CONFIGS_TOKEN"); token != "" {
		c.Token = token
	}
	if addr := os.Getenv("CONFIGS_ADDR"); addr != "" {
		c.Address = addr
	}
	if dc := os.Getenv("CONFIGS_DATACENTER"); dc != "" {
		c.DataCenter = dc
	}
}

func (c *configs) Load() error {
	c.mergeENV()
	config := api.DefaultConfig()
	config.Address = c.Address
	config.Token = c.Token
	config.Datacenter = c.DataCenter
	v, err := api.NewClient(config)
	if err != nil {
		return err
	}
	c.client = v
	return nil
}

func (c *configs) Prefix() string {
	return c.prefix
}

// Get returns the value of the key from consul
// WARNING: Please try to use SGet() to get the configs from consul
func (c *configs) Get(key string) ([]byte, error) {
	if c.client == nil {
		return nil, ErrNotInitialized
	}
	data, _, err := c.client.KV().Get(
		fmt.Sprintf("%s/%s", c.prefix, key),
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
func (c *configs) Set(key, val string) error {
	if c.client == nil {
		return ErrNotInitialized
	}
	_, err := c.client.KV().Put(
		&api.KVPair{
			Key:   fmt.Sprintf("%s/%s", c.prefix, key),
			Value: []byte(val),
		}, nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// SGet is recommended to use for getting configs from consul instead of using `Get()` directly.
// The main difference between `Get()` and `SGet()` is that the `SGet()` makes a limitation for the config's path.
// Before using this function to retrive the configs, you should make sure that your configs must be created at:
//	[environment]/services/[server name]/[your key]
// Ex. the service name is `userservice`.
//	   the config key is `abcd`
//	development/services/userservice/abcd
func (c *configs) SGet(k string) ([]byte, error) {
	if c.client == nil {
		return nil, ErrNotInitialized
	}
	data, _, err := c.client.KV().Get(
		fmt.Sprintf("%s/services/%s/%s", c.prefix, serverName, k),
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

// Client returns the consul client
func Client() *api.Client {
	return cfg.client
}

func Get() Configs {
	return &cfg
}
