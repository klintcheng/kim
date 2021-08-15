package gateway

import (
	"context"
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
	cmd.PersistentFlags().StringVarP(&opts.config, "config", "c", "conf.yaml", "Config file")
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
		Level: "info",
	})

	handler := &serv.Handler{
		ServiceID: config.ServiceID,
		AppSecret: config.AppSecret,
	}

	var srv kim.Server
	service := &naming.DefaultService{
		Id:       config.ServiceID,
		Name:     config.ServiceName,
		Address:  config.PublicAddress,
		Port:     config.PublicPort,
		Protocol: opts.protocol,
		Tags:     config.Tags,
	}

	if opts.protocol == "ws" {
		srv = websocket.NewServer(config.Listen, service)
	} else if opts.protocol == "tcp" {
		srv = tcp.NewServer(config.Listen, service)
	}

	srv.SetReadWait(time.Minute * 2)
	srv.SetAcceptor(handler)
	srv.SetMessageListener(handler)
	srv.SetStateListener(handler)

	_ = container.Init(srv, wire.SNChat, wire.SNLogin)

	ns, err := consul.NewNaming(config.ConsulURL)
	if err != nil {
		return err
	}
	container.SetServiceNaming(ns)

	// set a dialer
	container.SetDialer(serv.NewDialer(config.ServiceID))

	return container.Start()
}
