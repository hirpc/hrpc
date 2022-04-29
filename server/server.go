package server

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hirpc/hrpc/configs"
	"github.com/hirpc/hrpc/option"
	"github.com/hirpc/hrpc/utils/uniqueid"
	"google.golang.org/grpc"
)

var (
	// env is the current environment
	env = option.Development
	// name is the name of server
	name = ""
)

// Environment returns the current environment that the server is running for
func Environment() option.Environment {
	return env
}

// Name returns the name of this server
func Name() string {
	return name
}

// Server the server
type Server interface {
	// Server returns the gRPC server for registration at PB
	Server() *grpc.Server
	// Serve makes connections to databases and starts to listen ports to serve
	// it will block the current thread
	Serve() error
	// Run same like Serve() but it does not register to consul and starts servers to block
	// It will make connections to databases
	Run() error
}

func registeration(id, name string, port int, tags []string, healthCheck bool) error {
	reg := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Port:    port,
		Tags:    tags,
		Address: uniqueid.IP(),
	}

	if healthCheck {
		reg.Check = &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:6688", uniqueid.IP()),
			Timeout:                        "3s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "30s",
		}
		// for health check
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			runHealthAccepter()
		}()
	}
	return configs.Client().Agent().ServiceRegister(reg)
}

func deregisteration(id string) error {
	return configs.Client().Agent().ServiceDeregister(id)
}
