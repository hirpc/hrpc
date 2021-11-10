package server

import (
	"fmt"
	"net"

	"github.com/hirpc/hrpc/configs"
	"github.com/hirpc/hrpc/life"
	"github.com/hirpc/hrpc/option"
	"google.golang.org/grpc"
)

type GRPC struct {
	server *grpc.Server
	opts   *option.Options
}

func (g GRPC) Serve() error {
	if err := g.makeDatabase(); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", g.opts.ListenPort))
	if err != nil {
		return err
	}
	// registeration
	if err := registeration(
		g.opts.ID, g.opts.ServerName, g.opts.ListenPort,
		[]string{g.opts.ENV.String()},
		g.opts.HealthCheck,
	); err != nil {
		return err
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		g.server.Serve(lis)
	}()

	g.opts.WhenExit = append(g.opts.WhenExit, func() {
		deregisteration(g.opts.ID)
		g.server.Stop()
	})
	life.WhenExit(g.opts.WhenExit...)
	life.WhenRestart(g.opts.WhenRestart...)
	life.Start()
	return nil
}

func (g GRPC) Server() *grpc.Server {
	return g.server
}

func NewGRPC(opt *option.Options) (Server, error) {
	// register configs center
	if err := configs.Register(configs.Option{
		Token:      opt.ConsulCenter.Token,
		Address:    opt.ConsulCenter.Address,
		Prefix:     opt.ENV.String(),
		DataCenter: opt.ConsulCenter.DataCenter,
	}); err != nil {
		return nil, err
	}
	return &GRPC{
		server: grpc.NewServer(grpc.UnaryInterceptor(nil)),
		opts:   opt,
	}, nil
}
