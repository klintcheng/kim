package mock

import (
	"context"
	"fmt"
	"strings"

	"github.com/segmentio/ksuid"
	"github.com/spf13/cobra"
)

// StartOptions StartOptions
type StartOptions struct {
	addr     string
	protocol string
}

// NewCmd NewCmd
func NewClientCmd(ctx context.Context) *cobra.Command {
	opts := &StartOptions{}

	cmd := &cobra.Command{
		Use:   "mock_cli",
		Short: "start client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runcli(ctx, opts)
		},
	}
	cmd.PersistentFlags().StringVarP(&opts.addr, "address", "a", "localhost:8000", "server address")
	cmd.PersistentFlags().StringVarP(&opts.protocol, "protocol", "p", "ws", "protocol ws or tcp")
	return cmd
}

func runcli(ctx context.Context, opts *StartOptions) error {
	cli := ClientDemo{}
	if opts.protocol == "ws" && !strings.HasPrefix(opts.addr, "ws:") {
		opts.addr = fmt.Sprintf("ws://%s", opts.addr)
	}
	cli.Start(ksuid.New().String(), opts.protocol, opts.addr)
	return nil
}

// NewCmd NewCmd
func NewServerCmd(ctx context.Context) *cobra.Command {
	opts := &StartOptions{}

	cmd := &cobra.Command{
		Use:   "mock_srv",
		Short: "start server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runsrv(ctx, opts)
		},
	}
	cmd.PersistentFlags().StringVarP(&opts.addr, "address", "a", ":8000", "listen address")
	cmd.PersistentFlags().StringVarP(&opts.protocol, "protocol", "p", "ws", "protocol ws or tcp")
	return cmd
}

func runsrv(ctx context.Context, opts *StartOptions) error {
	srv := ServerDemo{}
	srv.Start("srv1", opts.protocol, opts.addr)
	return nil
}
