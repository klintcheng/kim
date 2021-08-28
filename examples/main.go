package main

import (
	"context"
	"flag"

	"github.com/klintcheng/kim/examples/echo"
	"github.com/klintcheng/kim/examples/mock"
	"github.com/klintcheng/kim/examples/throughput"
	"github.com/klintcheng/kim/logger"
	"github.com/spf13/cobra"
)

const version = "v1"

func main() {
	flag.Parse()

	root := &cobra.Command{
		Use:     "kim",
		Version: version,
		Short:   "server",
	}
	ctx := context.Background()

	// run echo client
	root.AddCommand(echo.NewCmd(ctx))

	// mock
	root.AddCommand(mock.NewClientCmd(ctx))
	root.AddCommand(mock.NewServerCmd(ctx))
	root.AddCommand(throughput.NewBenchmarkCmd(ctx))

	if err := root.Execute(); err != nil {
		logger.WithError(err).Fatal("Could not run command")
	}
}
