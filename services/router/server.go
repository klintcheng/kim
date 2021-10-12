package router

import (
	"context"
	"path"

	"github.com/kataras/iris/v12"
	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/naming/consul"
	"github.com/klintcheng/kim/services/router/apis"
	"github.com/klintcheng/kim/services/router/config"
	"github.com/klintcheng/kim/services/router/ipregion"
	"github.com/klintcheng/kim/services/service/conf"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ServerStartOptions ServerStartOptions
type ServerStartOptions struct {
	Listen    string
	ConsulURL string
	Path      string
}

// NewServerStartCmd creates a new http server command
func NewServerStartCmd(ctx context.Context, version string) *cobra.Command {
	opts := &ServerStartOptions{}

	cmd := &cobra.Command{
		Use:   "router",
		Short: "Start a router",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServerStart(ctx, opts, version)
		},
	}
	cmd.PersistentFlags().StringVarP(&opts.Listen, "listen", "l", ":8100", "listen hostPort")
	cmd.PersistentFlags().StringVarP(&opts.ConsulURL, "consul", "c", "localhost:8500", "consul url")
	cmd.PersistentFlags().StringVarP(&opts.Path, "path", "p", "./router", "base path")
	return cmd
}

// RunServerStart run http server
func RunServerStart(ctx context.Context, opts *ServerStartOptions, version string) error {
	_ = logger.Init(logger.Settings{
		Level: "trace",
	})

	ac := conf.MakeAccessLog()
	defer ac.Close()

	mappings, err := config.LoadMapping(path.Join(opts.Path, "mapping.json"))
	if err != nil {
		return err
	}
	logrus.Infof("load mappings - %v", mappings)
	regions, err := config.LoadRegions(path.Join(opts.Path, "regions.json"))
	if err != nil {
		return err
	}
	logrus.Infof("load regions - %v", regions)

	region, err := ipregion.NewIp2region(path.Join(opts.Path, "ip2region.db"))
	if err != nil {
		return err
	}

	ns, err := consul.NewNaming(opts.ConsulURL)
	if err != nil {
		return err
	}

	router := apis.RouterApi{
		Naming:   ns,
		IpRegion: region,
		Config: config.Router{
			Mapping: mappings,
			Regions: regions,
		},
	}

	app := iris.Default()
	app.UseRouter(ac.Handler)
	// app.UseRouter(setAllowedResponses)

	app.Get("/health", func(ctx iris.Context) {
		_, _ = ctx.WriteString("ok")
	})
	routerAPI := app.Party("/api/lookup")
	{
		routerAPI.Get("/:token", router.Lookup)
	}

	// Start server
	return app.Listen(opts.Listen, iris.WithOptimizations)
}

func setAllowedResponses(ctx iris.Context) {
	// Indicate that the Server can send JSON, XML, YAML and MessagePack for this request.
	ctx.Negotiation().JSON().MsgPack()
	// Add more, allowed by the server format of responses, mime types here...

	// If client is missing an "Accept: " header then default it to JSON.
	ctx.Negotiation().Accept.JSON()

	ctx.Next()
}
