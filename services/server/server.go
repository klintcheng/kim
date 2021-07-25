package server

import (
	"context"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/container"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/naming"
	"github.com/klintcheng/kim/naming/consul"
	"github.com/klintcheng/kim/services/server/conf"
	"github.com/klintcheng/kim/services/server/handler"
	"github.com/klintcheng/kim/services/server/serv"
	"github.com/klintcheng/kim/storage"
	"github.com/klintcheng/kim/tcp"
	"github.com/klintcheng/kim/wire"
	"github.com/spf13/cobra"
)

// ServerStartOptions ServerStartOptions
type ServerStartOptions struct {
	config      string
	serviceName string
}

// NewServerStartCmd creates a new http server command
func NewServerStartCmd(ctx context.Context, version string) *cobra.Command {
	opts := &ServerStartOptions{}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start a server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServerStart(ctx, opts, version)
		},
	}
	cmd.PersistentFlags().StringVarP(&opts.config, "config", "c", "./server/conf.yaml", "Config file")
	cmd.PersistentFlags().StringVarP(&opts.serviceName, "serviceName", "s", "chat", "defined a service name,option is login or chat")
	return cmd
}

// RunServerStart run http server
func RunServerStart(ctx context.Context, opts *ServerStartOptions, version string) error {
	config, err := conf.Init(opts.config)
	if err != nil {
		return err
	}
	_ = logger.Init(logger.Settings{
		Level: "trace",
	})

	r := kim.NewRouter()
	// login
	loginHandler := handler.NewLoginHandler()
	r.Handle(wire.CommandLoginSignIn, loginHandler.DoSysLogin)
	r.Handle(wire.CommandLoginSignOut, loginHandler.DoSysLogout)

	rdb, err := conf.InitRedis(config.RedisAddrs, "")
	if err != nil {
		return err
	}
	cache := storage.NewRedisStorage(rdb)
	servhandler := serv.NewServHandler(r, cache)

	service := &naming.DefaultService{
		Id:       config.ServiceID,
		Name:     opts.serviceName,
		Address:  config.PublicAddress,
		Port:     config.PublicPort,
		Protocol: string(wire.ProtocolTCP),
		Tags:     config.Tags,
	}
	srv := tcp.NewServer(config.Listen, service)

	srv.SetReadWait(kim.DefaultReadWait)
	srv.SetAcceptor(servhandler)
	srv.SetMessageListener(servhandler)
	srv.SetStateListener(servhandler)

	if err := container.Init(srv); err != nil {
		return err
	}

	ns, err := consul.NewNaming(config.ConsulURL)
	if err != nil {
		return err
	}
	container.SetServiceNaming(ns)

	return container.Start()
}
