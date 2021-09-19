package service

import (
	"context"
	"fmt"
	"hash/crc32"

	"gorm.io/gorm"

	"github.com/kataras/iris/v12"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/naming"
	"github.com/klintcheng/kim/naming/consul"
	"github.com/klintcheng/kim/services/service/conf"
	"github.com/klintcheng/kim/services/service/database"
	"github.com/klintcheng/kim/services/service/handler"
	"github.com/klintcheng/kim/wire"
	"github.com/spf13/cobra"
)

// ServerStartOptions ServerStartOptions
type ServerStartOptions struct {
	config string
}

// NewServerStartCmd creates a new http server command
func NewServerStartCmd(ctx context.Context, version string) *cobra.Command {
	opts := &ServerStartOptions{}

	cmd := &cobra.Command{
		Use:   "royal",
		Short: "Start a rpc service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServerStart(ctx, opts, version)
		},
	}
	cmd.PersistentFlags().StringVarP(&opts.config, "config", "c", "conf.yaml", "Config file")
	return cmd
}

// RunServerStart run http server
func RunServerStart(ctx context.Context, opts *ServerStartOptions, version string) error {
	config, err := conf.Init(opts.config)
	if err != nil {
		return err
	}
	_ = logger.Init(logger.Settings{
		Level:    config.LogLevel,
		Filename: "./data/royal.log",
	})

	// database.Init
	var (
		baseDb    *gorm.DB
		messageDb *gorm.DB
	)
	baseDb, err = database.InitDb(config.Driver, config.BaseDb)
	if err != nil {
		return err
	}
	messageDb, err = database.InitDb(config.Driver, config.MessageDb)
	if err != nil {
		return err
	}

	_ = baseDb.AutoMigrate(&database.Group{}, &database.GroupMember{})
	_ = messageDb.AutoMigrate(&database.MessageIndex{}, &database.MessageContent{})

	if config.NodeID == 0 {
		config.NodeID = int64(HashCode(config.ServiceID))
	}
	idgen, err := database.NewIDGenerator(config.NodeID)
	if err != nil {
		return err
	}

	rdb, err := conf.InitRedis(config.RedisAddrs, "")
	if err != nil {
		return err
	}

	ns, err := consul.NewNaming(config.ConsulURL)
	if err != nil {
		return err
	}
	_ = ns.Register(&naming.DefaultService{
		Id:       config.ServiceID,
		Name:     wire.SNService, // service name
		Address:  config.PublicAddress,
		Port:     config.PublicPort,
		Protocol: "http",
		Tags:     config.Tags,
		Meta: map[string]string{
			consul.KeyHealthURL: fmt.Sprintf("http://%s:%d/health", config.PublicAddress, config.PublicPort),
		},
	})
	defer func() {
		_ = ns.Deregister(config.ServiceID)
	}()
	serviceHandler := handler.ServiceHandler{
		BaseDb:    baseDb,
		MessageDb: messageDb,
		Idgen:     idgen,
		Cache:     rdb,
	}

	ac := conf.MakeAccessLog()
	defer ac.Close()

	app := newApp(&serviceHandler)
	app.UseRouter(ac.Handler)
	app.UseRouter(setAllowedResponses)

	// Start server
	return app.Listen(config.Listen, iris.WithOptimizations)
}

func newApp(serviceHandler *handler.ServiceHandler) *iris.Application {
	app := iris.Default()

	app.Get("/health", func(ctx iris.Context) {
		_, _ = ctx.WriteString("ok")
	})
	messageAPI := app.Party("/api/:app/message")
	{
		messageAPI.Post("/user", serviceHandler.InsertUserMessage)
		messageAPI.Post("/group", serviceHandler.InsertGroupMessage)
		messageAPI.Post("/ack", serviceHandler.MessageAck)
	}

	groupAPI := app.Party("/api/:app/group")
	{
		groupAPI.Get("/:id", serviceHandler.GroupGet)
		groupAPI.Post("", serviceHandler.GroupCreate)
		groupAPI.Post("/member", serviceHandler.GroupJoin)
		groupAPI.Delete("/member", serviceHandler.GroupQuit)
		groupAPI.Get("/members/:id", serviceHandler.GroupMembers)
	}

	offlineAPI := app.Party("/api/:app/offline")
	{
		offlineAPI.Use(iris.Compression)
		offlineAPI.Post("/index", serviceHandler.GetOfflineMessageIndex)
		offlineAPI.Post("/content", serviceHandler.GetOfflineMessageContent)
	}
	return app
}

func setAllowedResponses(ctx iris.Context) {
	// Indicate that the Server can send JSON, XML, YAML and MessagePack for this request.
	ctx.Negotiation().JSON().Protobuf().MsgPack()
	// Add more, allowed by the server format of responses, mime types here...

	// If client is missing an "Accept: " header then default it to JSON.
	ctx.Negotiation().Accept.JSON()

	ctx.Next()
}

func HashCode(key string) uint32 {
	hash32 := crc32.NewIEEE()
	hash32.Write([]byte(key))
	return hash32.Sum32() % 1000
}
