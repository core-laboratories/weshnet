package main

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/peterbourgon/ff/v3/ffcli"
	"go.uber.org/zap"

	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/errcode"
	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/messengertypes"
)

func remoteLogsCommand() *ffcli.Command {
	fsBuilder := func() (*flag.FlagSet, error) {
		fs := flag.NewFlagSet("wesh remote-logs", flag.ExitOnError)
		manager.SetupRemoteNodeFlags(fs)
		return fs, nil
	}

	return &ffcli.Command{
		Name:           "remote-logs",
		ShortUsage:     "wesh [global flags] remote-logs",
		ShortHelp:      "stream logs from a remote node",
		Options:        ffSubcommandOptions(),
		FlagSetBuilder: fsBuilder,
		UsageFunc:      usageFunc,
		Exec: func(ctx context.Context, args []string) error { //nolint:revive
			manager.SetLogger(zap.NewNop())
			manager.DisableIPFSNetwork()

			// messenger
			messenger, err := manager.GetMessengerClient()
			if err != nil {
				return err
			}

			stream, err := messenger.DevStreamLogs(ctx, &messengertypes.DevStreamLogs_Request{})
			if err != nil {
				return err
			}

			for {
				ret, err := stream.Recv()
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return errcode.ErrCode_TODO.Wrap(err)
				}
				fmt.Println(ret.Line)
			}
		},
	}
}
