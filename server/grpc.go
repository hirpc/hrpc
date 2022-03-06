package server

import (
	"crypto/tls"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/hirpc/hrpc/life"
	"github.com/hirpc/hrpc/log"
	"github.com/hirpc/hrpc/option"
	"github.com/hirpc/hrpc/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type HRPC struct {
	server *grpc.Server
	opts   *option.Options
}

func (h HRPC) Run() error {
	return h.makeDatabase()
}

func (h HRPC) Serve() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if err := h.makeDatabase(); err != nil {
		return err
	}
	if err := h.makeMessageQueue(); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", h.opts.ListenPort))
	if err != nil {
		return err
	}
	// registeration
	if err := registeration(
		h.opts.ID, h.opts.ServerName, h.opts.ListenPort,
		[]string{h.opts.ENV.String()},
		h.opts.HealthCheck,
	); err != nil {
		return err
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		h.server.Serve(lis)
	}()

	h.opts.WhenExit = append(h.opts.WhenExit, func() {
		deregisteration(h.opts.ID)
		h.server.Stop()
	})
	h.opts.WhenRestart = append(h.opts.WhenRestart, func() {
		deregisteration(h.opts.ID)
		h.server.Stop()
	})
	life.WhenExit(h.opts.WhenExit...)
	life.WhenRestart(h.opts.WhenRestart...)
	life.Start()
	return nil
}

func (h HRPC) Server() *grpc.Server {
	return h.server
}

func grpcOption(opt *option.Options) ([]grpc.ServerOption, error) {
	gopt := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				tracer.AddTraceID,
				log.AuditLog,
			),
		),
	}
	if opt.ServerCerts != nil {
		cert, err := tls.X509KeyPair(opt.ServerCerts.PubKey, opt.ServerCerts.PriKey)
		if err != nil {
			return nil, err
		}
		gopt = append(gopt, grpc.Creds(credentials.NewServerTLSFromCert(&cert)))
	}
	return gopt, nil
}

func NewHRPC(opt *option.Options) (Server, error) {
	// set current environment
	env = opt.ENV
	name = opt.ServerName
	gopt, err := grpcOption(opt)
	if err != nil {
		return nil, err
	}
	return &HRPC{
		server: grpc.NewServer(gopt...),
		opts:   opt,
	}, nil
}
