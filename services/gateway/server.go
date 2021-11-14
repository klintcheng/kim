package gateway

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/container"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/naming"
	"github.com/klintcheng/kim/naming/consul"
	"github.com/klintcheng/kim/services/gateway/conf"
	"github.com/klintcheng/kim/services/gateway/serv"
	"github.com/klintcheng/kim/tcp"
	"github.com/klintcheng/kim/websocket"
	"github.com/klintcheng/kim/wire"
	"github.com/spf13/cobra"
)

// const logName = "logs/gateway"

// ServerStartOptions ServerStartOptions
type ServerStartOptions struct {
	config   string
	protocol string
	route    string
}

// NewServerStartCmd creates a new http server command
func NewServerStartCmd(ctx context.Context, version string) *cobra.Command {
	opts := &ServerStartOptions{}

	cmd := &cobra.Command{
		Use:   "gateway",
		Short: "Start a gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServerStart(ctx, opts, version)
		},
	}
	cmd.PersistentFlags().StringVarP(&opts.config, "config", "c", "./gateway/conf.yaml", "Config file")
	cmd.PersistentFlags().StringVarP(&opts.route, "route", "r", "./gateway/route.json", "route file")
	cmd.PersistentFlags().StringVarP(&opts.protocol, "protocol", "p", "ws", "protocol of ws or tcp")
	return cmd
}

// RunServerStart run http server
func RunServerStart(ctx context.Context, opts *ServerStartOptions, version string) error {
	config, err := conf.Init(opts.config)
	if err != nil {
		return err
	}
	_ = logger.Init(logger.Settings{
		Level:    "trace",
		Filename: "./data/gateway.log",
	})

	handler := &serv.Handler{
		ServiceID: config.ServiceID,
		AppSecret: config.AppSecret,
	}
	meta := make(map[string]string)
	meta[consul.KeyHealthURL] = fmt.Sprintf("http://%s:%d/health", config.PublicAddress, config.MonitorPort)
	meta["domain"] = config.Domain

	var srv kim.Server
	service := &naming.DefaultService{
		Id:       config.ServiceID,
		Name:     config.ServiceName,
		Address:  config.PublicAddress,
		Port:     config.PublicPort,
		Protocol: opts.protocol,
		Tags:     config.Tags,
		Meta:     meta,
	}
	srvOpts := []kim.ServerOption{
		kim.WithConnectionGPool(config.ConnectionGPool), kim.WithMessageGPool(config.MessageGPool),
	}
	if opts.protocol == "ws" {
		srv = websocket.NewServer(config.Listen, service, srvOpts...)
	} else if opts.protocol == "tcp" {
		srv = tcp.NewServer(config.Listen, service, srvOpts...)
	}

	srv.SetReadWait(time.Minute * 2)
	srv.SetAcceptor(handler)
	srv.SetMessageListener(handler)
	srv.SetStateListener(handler)

	_ = container.Init(srv, wire.SNChat, wire.SNLogin)
	container.EnableMonitor(fmt.Sprintf(":%d", config.MonitorPort))

	ns, err := consul.NewNaming(config.ConsulURL)
	if err != nil {
		return err
	}
	container.SetServiceNaming(ns)
	// set a dialer
	container.SetDialer(serv.NewDialer(config.ServiceID))
	// use routeSelector
	selector, err := serv.NewRouteSelector(opts.route)
	if err != nil {
		return err
	}
	container.SetSelector(selector)
	return container.Start()
}
