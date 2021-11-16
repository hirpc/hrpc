package option

import "github.com/hirpc/hrpc/plugin"

func WithPlugins(plugins ...plugin.Plugin) Option {
	return func(o *Options) {
		o.Plugins = plugins
	}
}
