package service

import (
	"fmt"

	"github.com/hirpc/hrpc/configs"
	"github.com/hirpc/hrpc/option"
	"github.com/hirpc/hrpc/server"
)

// Service represents a service
type Service struct {
	ID   string
	Name string
	// Endpoint represents the IP address of this program is running
	Endpoint string
	Weight   int
	Port     int
}

// String to string for connection with POD IP:Port
func (s Service) String() string {
	return fmt.Sprintf("%s:%d", s.Endpoint, s.Port)
}

// Target returns a string that should be used in TKE environment only
// since we can use DNS feature to find out the service IP to viste in case of the changes of POD IP
// !!! The namespace in QCloud must be prod (stands for production) OR dev (stands for development).
// Otherwise, the request can be failed.
func (s Service) Target() string {
	var namespace = "prod"
	if server.Environment() == option.Development {
		namespace = "dev"
	}
	return fmt.Sprintf("%s.%s.svc.cluster.local:%d", s.Name, namespace, s.Port)
}

// Get will return the best service based on the name and tag
// If the tag provides more than one, it will only pick the first one
// If the tag does not provide, the service will be picked based on the current environment(prefix of configs)
func Get(name string, tags ...Tag) (*Service, error) {
	var tag Tag
	if len(tags) != 0 {
		tag = tags[0]
	} else {
		tag = Tag(configs.Get().Prefix())
	}
	ss, _, err := configs.Client().Health().Service(name, tag.String(), true, nil)
	if err != nil {
		return nil, err
	}
	if len(ss) == 0 {
		return nil, fmt.Errorf("[%s] not found or no permission to read", name)
	}

	var maxWeight int
	var index int
	for i, s := range ss {
		if s.Service.Weights.Passing > maxWeight {
			maxWeight = s.Service.Weights.Passing
			index = i
		}
	}
	return &Service{
		ID:       ss[index].Service.ID,
		Name:     ss[index].Service.Service,
		Endpoint: ss[index].Service.Address,
		Weight:   maxWeight,
		Port:     ss[index].Service.Port,
	}, nil
}
