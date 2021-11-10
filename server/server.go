package server

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hirpc/arsenal/uniqueid"
	"github.com/hirpc/hrpc/configs"
	"google.golang.org/grpc"
)

type Server interface {
	Serve() error
	Server() *grpc.Server
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
