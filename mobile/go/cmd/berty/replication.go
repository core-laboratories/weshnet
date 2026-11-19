package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	legacygateway "github.com/grpc-ecosystem/grpc-gateway/runtime"
	gateway "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/authtypes"
	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/bertyreplication"
	"github.com/core-laboratories/weshnet/v2/pkg/replicationtypes"
)

func replicationServerCommand() *ffcli.Command {
	fsBuilder := func() (*flag.FlagSet, error) {
		fs := flag.NewFlagSet("berty repl-server", flag.ExitOnError)
		fs.String("config", "", "config file (optional)")
		manager.Session.Kind = "cli.replication"
		manager.SetupLoggingFlags(fs) // also available at root level
		manager.SetupProtocolAuth(fs)
		manager.SetupLocalProtocolServerFlags(fs)
		manager.SetupDefaultGRPCListenersFlags(fs)

		// set serviceid for needed by push server
		manager.Node.Protocol.ServiceID = authtypes.ServiceReplicationID
		return fs, nil
	}

	return &ffcli.Command{
		Name:           "repl-server",
		ShortHelp:      "replication server",
		ShortUsage:     "berty [global flags] repl-server [flags]",
		FlagSetBuilder: fsBuilder,
		UsageFunc:      usageFunc,
		Options:        ffSubcommandOptions(),
		Exec: func(ctx context.Context, args []string) error {
			if len(args) > 0 {
				return flag.ErrHelp
			}

			if manager.Node.Protocol.AuthSecret == "" {
				return fmt.Errorf("node.auth-secret cannot be empty")
			}

			if manager.Node.Protocol.AuthPublicKey == "" {
				return fmt.Errorf("node.auth-pk cannot be empty")
			}

			var err error
			server, mux, err := manager.GetGRPCServer()
			if err != nil {
				return err
			}

			logger, err := manager.GetLogger()
			if err != nil {
				return err
			}

			db, err := manager.GetReplicationDB()
			if err != nil {
				return err
			}

			odb, err := manager.GetOrbitDB()
			if err != nil {
				return err
			}

			replicationService, err := bertyreplication.NewReplicationService(ctx, db, odb, logger)
			if err != nil {
				return err
			}

			replicationtypes.RegisterReplicationServiceServer(server, replicationService)
		replicationMux := gateway.NewServeMux()
		if err := replicationtypes.RegisterReplicationServiceHandlerServer(ctx, replicationMux, replicationService); err != nil {
				return err
			}
		for _, method := range []string{"ReplicateGroup", "ReplicateGlobalStats", "ReplicateGroupStats"} {
			pattern := legacygateway.MustPattern(legacygateway.NewPattern(1, []int{2, 0, 2, 1}, []string{"weshnet.replication.v1.ReplicationService", method}, ""))
			mux.Handle(http.MethodPost, pattern, func(w http.ResponseWriter, req *http.Request, _ map[string]string) {
				replicationMux.ServeHTTP(w, req)
			})
		}

			return manager.RunWorkers(ctx)
		},
	}
}
