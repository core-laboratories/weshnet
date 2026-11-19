package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/weshdirectory"
	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/directorytypes"
)

func directoryServiceCommand() *ffcli.Command {
	allowedIssuers := ""

	fsBuilder := func() (*flag.FlagSet, error) {
		fs := flag.NewFlagSet("wesh directory-service", flag.ExitOnError)
		manager.Session.Kind = "cli.directory-service"
		manager.SetupLoggingFlags(fs) // also available at root level
		manager.SetupProtocolAuth(fs)
		manager.SetupLocalProtocolServerFlags(fs)
		manager.SetupDefaultGRPCListenersFlags(fs)
		fs.StringVar(&allowedIssuers, "allowed-issuers", allowedIssuers, "Allowed verified credentials issuers, comma-separated")
		return fs, nil
	}

	return &ffcli.Command{
		Name:           "directory-service",
		ShortUsage:     "wesh [global flags] directory-service [flags]",
		ShortHelp:      "starts a directory service",
		Options:        ffSubcommandOptions(),
		FlagSetBuilder: fsBuilder,
		UsageFunc:      usageFunc,
		Exec: func(ctx context.Context, args []string) error { //nolint:revive
			if allowedIssuers == "" {
				return fmt.Errorf("allowed-issuers cannot be empty")
			}

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			server, mux, err := manager.GetGRPCServer()
			if err != nil {
				return err
			}

			logger, err := manager.GetLogger()
			if err != nil {
				return err
			}

			db, err := manager.GetDirectoryServiceDB()
			if err != nil {
				return err
			}

			dsConfig := &weshdirectory.ServiceOpts{
				Logger:         logger,
				AllowedIssuers: strings.Split(allowedIssuers, ","),
			}

			ds, err := weshdirectory.NewSQLDatastore(db)
			if err != nil {
				return err
			}

			directoryService, err := weshdirectory.New(ds, dsConfig)
			if err != nil {
				return err
			}

			directorytypes.RegisterDirectoryServiceServer(server, directoryService)
			if err := directorytypes.RegisterDirectoryServiceHandlerServer(ctx, mux, directoryService); err != nil {
				return err
			}

			return manager.RunWorkers(ctx)
		},
	}
}
