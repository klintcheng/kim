package throughput

import (
	"context"
	"runtime"

	"github.com/spf13/cobra"
)

// StartOptions StartOptions
type StartOptions struct {
	addr    string
	chat    string
	count   int
	offline bool
}

// NewCmd NewCmd
func NewBenchmarkCmd(ctx context.Context) *cobra.Command {
	opts := &StartOptions{}

	cmd := &cobra.Command{
		Use:   "benchmark",
		Short: "start client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runcli(ctx, opts)
		},
	}
	cmd.PersistentFlags().StringVarP(&opts.addr, "address", "a", "ws://localhost:8000", "server address")
	cmd.PersistentFlags().StringVarP(&opts.chat, "chattype", "c", "user", "user or group")
	cmd.PersistentFlags().IntVarP(&opts.count, "number", "n", 100, "message number")
	cmd.PersistentFlags().BoolVarP(&opts.offline, "offline", "o", true, "receiver offline")
	return cmd
}

func runcli(ctx context.Context, opts *StartOptions) error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if opts.chat == "user" {
		err := usertalk(opts.addr, opts.count, opts.offline)
		if err != nil {
			return err
		}
	}
	return nil
}
