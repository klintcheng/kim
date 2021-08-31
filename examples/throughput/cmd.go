package throughput

import (
	"context"
	"time"

	"github.com/klintcheng/kim/wire/token"
	"github.com/spf13/cobra"
)

// DefaultOptions DefaultOptions
type Options struct {
	Addr      string
	AppSecret string
	Count     int
}

// NewCmd NewCmd
func NewBenchmarkCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "benchmark",
		Short: "start client",
	}
	var opts = &Options{}
	cmd.PersistentFlags().StringVarP(&opts.Addr, "address", "a", "ws://124.71.204.19:8000", "server address")
	cmd.PersistentFlags().StringVarP(&opts.AppSecret, "appSecret", "s", token.DefaultSecret, "app secret")
	cmd.PersistentFlags().IntVarP(&opts.Count, "number", "n", 100, "message number")

	cmd.AddCommand(NewUserTalkCmd(opts))
	cmd.AddCommand(NewGroupTalkCmd(opts))
	cmd.AddCommand(NewLoginCmd(opts))
	return cmd
}

type UserOptions struct {
	Offline bool
}

func NewUserTalkCmd(opts *Options) *cobra.Command {
	var options = &UserOptions{}

	cmd := &cobra.Command{
		Use:   "user",
		Short: "u",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := usertalk(opts.Addr, opts.AppSecret, opts.Count, options.Offline)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&options.Offline, "offline", "o", true, "receiver offline")
	return cmd
}

type LoginOptions struct {
	keep time.Duration
}

func NewLoginCmd(opts *Options) *cobra.Command {
	var options = &LoginOptions{}
	cmd := &cobra.Command{
		Use:   "login",
		Short: "lo",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := login(opts.Addr, opts.AppSecret, opts.Count, options.keep)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.PersistentFlags().DurationVarP(&options.keep, "keep", "k", time.Second*5, "the duration of keeping the client connection")
	return cmd
}

type GroupOptions struct {
	MemberCount   int
	OnlinePercent float32
}

func NewGroupTalkCmd(opts *Options) *cobra.Command {
	var options = &GroupOptions{}

	cmd := &cobra.Command{
		Use:   "group",
		Short: "gp",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := grouptalk(opts.Addr, opts.AppSecret, opts.Count, options.MemberCount, options.OnlinePercent)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.PersistentFlags().IntVarP(&options.MemberCount, "memcount", "m", 20, "member count")
	cmd.PersistentFlags().Float32VarP(&options.OnlinePercent, "percet", "p", 0.2, "online percet")
	return cmd
}
