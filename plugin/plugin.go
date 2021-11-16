package plugin

import (
	"sync"

	"github.com/hirpc/hrpc/uerror"
)

var (
	ErrMissingPlugin = uerror.New(1001, "dependence check failed because missing plugins")
	ErrDependsCycle  = uerror.New(1002, "depends cycle")
)

type Plugin interface {
	Load() error
	Name() string

	DependsOn() []string
}

type plugin struct {
	plugins        map[string]Plugin
	orderedPlugins chan Plugin
	statuses       map[string]bool

	mu sync.RWMutex
}

var p = plugin{
	plugins:  map[string]Plugin{},
	statuses: map[string]bool{},
}

func Register(plugins ...Plugin) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, plugin := range plugins {
		p.plugins[plugin.Name()] = plugin
		p.statuses[plugin.Name()] = false
	}
}

func Setup() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.orderedPlugins = make(chan Plugin, len(p.plugins))
	for _, plugin := range p.plugins {
		p.orderedPlugins <- plugin
	}

	num := len(p.orderedPlugins)
	for num > 0 {
		for i := 0; i < num; i++ {
			plug := <-p.orderedPlugins
			if plug.DependsOn() == nil || len(plug.DependsOn()) == 0 {
				if err := plug.Load(); err != nil {
					return err
				}
				p.statuses[plug.Name()] = true
			} else {
				allowed := true
				for _, dependName := range plug.DependsOn() {
					if !p.statuses[dependName] {
						// has depends that unloaded
						allowed = false
						p.orderedPlugins <- plug
						break
					}
				}
				if allowed {
					if err := plug.Load(); err != nil {
						return err
					}
					p.statuses[plug.Name()] = true
				}
			}
		}
		if len(p.orderedPlugins) == num {
			return ErrDependsCycle
		}
		num = len(p.orderedPlugins)
	}
	return nil
}
