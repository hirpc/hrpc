package hrpc

import (
	"context"
	"time"

	"github.com/hirpc/arsenal/uniqueid"
	"github.com/hirpc/hrpc/codec"
	"github.com/hirpc/hrpc/configs"
	"github.com/hirpc/hrpc/database"
	"github.com/hirpc/hrpc/log"
	"github.com/hirpc/hrpc/log/hook"
	"github.com/hirpc/hrpc/option"
	"github.com/hirpc/hrpc/plugin"
	"github.com/hirpc/hrpc/server"
	"github.com/hirpc/hrpc/tracer"
	"github.com/hirpc/hrpc/utils/location"
	"google.golang.org/grpc/metadata"
)

func load() error {
	if err := uniqueid.Register(); err != nil {
		return err
	}
	return nil
}

// NewServer is the entrance of the framwork
func NewServer(opts ...option.Option) (server.Server, error) {
	// register some dependent components
	if err := load(); err != nil {
		return nil, err
	}

	var opt = &option.Options{
		ID:          uniqueid.String(),
		ListenPort:  8888,
		ENV:         option.Development,
		DBs:         make(map[string]database.Database),
		HealthCheck: false,
		StackSkip:   1,
	}
	for _, o := range opts {
		o(opt)
	}
	if err := opt.Valid(); err != nil {
		return nil, err
	}

	plugin.Register(location.New(), configs.New(
		opt.ConsulCenter.Address,
		opt.ConsulCenter.DataCenter,
		opt.ENV.String(),
		opt.ConsulCenter.Token,
	))
	plugin.Register(opt.Plugins...)
	if err := plugin.Setup(); err != nil {
		return nil, err
	}

	if hook.CLSHook() != nil {
		// init log component
		log.With(log.Option{
			Environment: opt.ENV.String(),
			Hooks:       []hook.Hook{hook.CLSHook()},
			StackSkip:   opt.StackSkip,
		})
	} else {
		// init log component
		log.With(log.Option{
			Environment: opt.ENV.String(),
			StackSkip:   opt.StackSkip,
		})
	}

	// returns a new grpc server with desired options
	return server.NewHRPC(opt)
}

// BackgroundContext should be used to instead context.Context
func BackgroundContext() context.Context {
	ctx := context.Background()
	msg := codec.Message(ctx)
	msg.WithTraceID(tracer.NewID())
	msg.WithServerName(server.Name())
	msg.WithNamespace(server.Environment().String())
	msg.WithRequestTimeout(time.Second * 3)
	ctx = metadata.NewOutgoingContext(
		ctx, msg.Metadata(),
	)
	return metadata.NewIncomingContext(
		ctx, msg.Metadata(),
	)
}
