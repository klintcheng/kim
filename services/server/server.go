package server

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/container"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/naming"
	"github.com/klintcheng/kim/naming/consul"
	"github.com/klintcheng/kim/services/server/conf"
	"github.com/klintcheng/kim/services/server/handler"
	"github.com/klintcheng/kim/services/server/serv"
	"github.com/klintcheng/kim/services/server/service"
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
	cmd.PersistentFlags().StringVarP(&opts.config, "config", "c", "conf.yaml", "Config file")
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
		Level:    "info",
		Filename: "./data/server.log",
	})

	var groupService service.Group
	var messageService service.Message
	if config.RoyalURL != "" {
		groupService = service.NewGroupService(config.RoyalURL)
		messageService = service.NewMessageService(config.RoyalURL)
	} else {
		srvRecord := &resty.SRVRecord{
			Domain:  "consul",
			Service: wire.SNService,
		}
		groupService = service.NewGroupServiceWithSRV("http", srvRecord)
		messageService = service.NewMessageServiceWithSRV("http", srvRecord)
	}

	r := kim.NewRouter()
	// login
	loginHandler := handler.NewLoginHandler()
	r.Handle(wire.CommandLoginSignIn, loginHandler.DoSysLogin)
	r.Handle(wire.CommandLoginSignOut, loginHandler.DoSysLogout)
	// talk
	chatHandler := handler.NewChatHandler(messageService, groupService)
	r.Handle(wire.CommandChatUserTalk, chatHandler.DoUserTalk)
	r.Handle(wire.CommandChatGroupTalk, chatHandler.DoGroupTalk)
	r.Handle(wire.CommandChatTalkAck, chatHandler.DoTalkAck)
	// group
	groupHandler := handler.NewGroupHandler(groupService)
	r.Handle(wire.CommandGroupCreate, groupHandler.DoCreate)
	r.Handle(wire.CommandGroupJoin, groupHandler.DoJoin)
	r.Handle(wire.CommandGroupQuit, groupHandler.DoQuit)
	r.Handle(wire.CommandGroupDetail, groupHandler.DoDetail)

	// offline
	offlineHandler := handler.NewOfflineHandler(messageService)
	r.Handle(wire.CommandOfflineIndex, offlineHandler.DoSyncIndex)
	r.Handle(wire.CommandOfflineContent, offlineHandler.DoSyncContent)

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
