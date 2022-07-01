package hrpc

import (
	"context"
	"time"

	"github.com/hirpc/hrpc/codec"
	"github.com/hirpc/hrpc/configs"
	"github.com/hirpc/hrpc/database"
	"github.com/hirpc/hrpc/log"
	"github.com/hirpc/hrpc/log/cls"
	"github.com/hirpc/hrpc/mq"
	"github.com/hirpc/hrpc/option"
	"github.com/hirpc/hrpc/plugin"
	"github.com/hirpc/hrpc/server"
	"github.com/hirpc/hrpc/tracer"
	"github.com/hirpc/hrpc/utils/location"
	"github.com/hirpc/hrpc/utils/uniqueid"
	"google.golang.org/grpc/metadata"
)

// NewServer is the entrance of the framwork
func NewServer(opts ...option.Option) (server.Server, error) {
	var opt = &option.Options{
		ListenPort:  8888,
		ENV:         option.Development,
		DBs:         make(map[string]database.Database),
		MQs:         make(map[string]mq.MQ),
		HealthCheck: false,
		StackSkip:   1,
	}
	for _, o := range opts {
		o(opt)
	}
	if err := opt.Valid(); err != nil {
		return nil, err
	}

	plugin.Register(uniqueid.New(), location.New(), configs.New(
		opt.ConsulCenter.Address,
		opt.ConsulCenter.DataCenter,
		opt.ENV.String(),
		opt.ConsulCenter.Token,
		opt.ServerName,
	))
	plugin.Register(opt.Plugins...)
	if err := plugin.Setup(); err != nil {
		return nil, err
	}

	// fixup the ID
	opt.ID = uniqueid.String()

	logOpt := log.Option{
		Environment: opt.ENV.String(),
		StackSkip:   opt.StackSkip,
	}
	if cls.Hook() != nil {
		logOpt.Hooks = append(logOpt.Hooks, cls.Hook())
	}
	log.With(logOpt)

	// returns a new grpc server with desired options
	return server.NewHRPC(opt)
}

// BackgroundContext should be used to instead context.Context
func BackgroundContext() context.Context {
	ctx := context.Background()
	msg := codec.Message(ctx)
	msg.WithServerName(server.Name())
	msg.WithNamespace(server.Environment().String())
	msg.WithRequestTimeout(time.Second * 3)
	msg.WithTraceID(tracer.NewID(server.Name()))
	ctx = metadata.NewOutgoingContext(
		ctx, msg.Metadata(),
	)
	return metadata.NewIncomingContext(
		ctx, msg.Metadata(),
	)
}
