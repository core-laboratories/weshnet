package main

import (
	"context"
	"flag"
	"fmt"
	"runtime"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/weshversion"
)

func versionCommand() *ffcli.Command {
	return &ffcli.Command{
		Name:       "version",
		ShortUsage: "wesh version",
		ShortHelp:  "print software version",
		Options:    ffSubcommandOptions(),
		UsageFunc:  usageFunc,
		Exec: func(_ context.Context, args []string) error {
			if len(args) > 0 {
				return flag.ErrHelp
			}
			fmt.Printf("version  %s\n", weshversion.Version)
			if weshversion.VcsRef != "n/a" {
				fmt.Printf("vcs      https://github.com/core-laboratories/weshnet/commit/%s\n", weshversion.VcsRef)
			}
			fmt.Printf("go       %s\n", runtime.Version())
			return nil
		},
	}
}
