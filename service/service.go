package service

import (
	"errors"
	"fmt"

	"github.com/hirpc/hrpc/configs"
)

// Service represents a service
type Service struct {
	ID       string
	Name     string
	Endpoint string
	Weight   int
	Port     int
}

// String to string for connection
func (s Service) String() string {
	return fmt.Sprintf("%s:%d", s.Endpoint, s.Port)
}

var (
	// ErrNotFound not found any services
	ErrNotFound = errors.New("not found")
)

// Get will return the best service based on the name and tag
// If the tag provides more than one, it will only pick the first one
// If the tag does not provide, the service will be picked based on the current environment(prefix of configs)
func Get(name string, tags ...Tag) (*Service, error) {
	var tag Tag
	if len(tags) != 0 {
		tag = tags[0]
	} else {
		tag = Tag(configs.Prefix())
	}
	ss, _, err := configs.Client().Health().Service(name, tag.String(), true, nil)
	if err != nil {
		return nil, err
	}
	if len(ss) == 0 {
		return nil, ErrNotFound
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

// List returns a set of services based on the its name
// If the tag does not provided, it will pick the configs.Prefix() valud which represents the current environment
func List(name string, tags ...Tag) ([]*Service, error) {
	var tag Tag
	if len(tags) != 0 {
		tag = tags[0]
	} else {
		tag = Tag(configs.Prefix())
	}
	ss, _, err := configs.Client().Health().Service(name, tag.String(), true, nil)
	if err != nil {
		return nil, err
	}
	if len(ss) == 0 {
		return nil, ErrNotFound
	}

	var services = make([]*Service, len(ss))
	for i, s := range ss {
		services[i] = &Service{
			ID:       s.Service.ID,
			Name:     s.Service.Service,
			Endpoint: s.Service.Address,
			Weight:   s.Service.Weights.Passing,
			Port:     s.Service.Port,
		}
	}
	return services, nil
}
