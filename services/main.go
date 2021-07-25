package main

import (
	"context"
	"flag"

	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/services/gateway"
	"github.com/klintcheng/kim/services/server"
	"github.com/spf13/cobra"
)

const version = "v1"

func main() {
	flag.Parse()

	root := &cobra.Command{
		Use:     "kim",
		Version: version,
		Short:   "King IM Cloud",
	}
	ctx := context.Background()

	root.AddCommand(gateway.NewServerStartCmd(ctx, version))
	root.AddCommand(server.NewServerStartCmd(ctx, version))

	if err := root.Execute(); err != nil {
		logger.WithError(err).Fatal("Could not run command")
	}
}
