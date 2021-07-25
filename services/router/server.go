package router

import (
	"context"
	"hash/crc32"

	"github.com/kataras/iris/v12"
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
		Use:   "router",
		Short: "Start a router",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServerStart(ctx, opts, version)
		},
	}
	cmd.PersistentFlags().StringVarP(&opts.config, "config", "c", "./router/conf.yaml", "Config file")
	return cmd
}

// RunServerStart run http server
func RunServerStart(ctx context.Context, opts *ServerStartOptions, version string) error {

	app := iris.Default()

	app.Get("/health", func(ctx iris.Context) {
		_, _ = ctx.WriteString("ok")
	})

	// Start server
	return app.Listen("", iris.WithOptimizations)
}

func HashCode(key string) uint32 {
	hash32 := crc32.NewIEEE()
	hash32.Write([]byte(key))
	return hash32.Sum32() % 1000
}
