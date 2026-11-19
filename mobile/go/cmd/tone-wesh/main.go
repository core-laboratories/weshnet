package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	"go.uber.org/zap"

	"github.com/core-laboratories/weshnet/mobile/v2/go/internal/grpcserver"
	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/accounttypes"
	accountsvc "github.com/core-laboratories/weshnet/mobile/v2/go/pkg/weshaccount"
)

func main() {
	var accountListeners string
	var serviceListeners string
	var appDir string
	var sharedDir string
	flag.StringVar(&accountListeners, "node.account.listeners", "/ip4/127.0.0.1/tcp/9092/grpc", "local account API listener")
	flag.StringVar(&serviceListeners, "node.listeners", "/ip4/127.0.0.1/tcp/9091/grpc", "local messenger API listener")
	flag.StringVar(&appDir, "store.dir", "", "account data directory")
	flag.StringVar(&sharedDir, "store.shared-dir", "", "shared account data directory")
	flag.String("log.filters", "", "accepted for compatibility; Tone keeps runtime logs private")
	flag.Parse()
	if appDir == "" || sharedDir == "" {
		fmt.Fprintln(os.Stderr, "store.dir and store.shared-dir are required")
		os.Exit(2)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	logger := zap.NewNop()
	workers := run.Group{}
	workers.Add(func() error {
		<-ctx.Done()
		return ctx.Err()
	}, func(error) {
		cancel()
	})

	server, gateway, _, err := grpcserver.InitGRPCServer(&workers, &grpcserver.GRPCOpts{
		Logger:    logger,
		Listeners: accountListeners,
	})
	if err != nil {
		fatal(err)
	}
	service, err := accountsvc.NewService(&accountsvc.Options{
		AppRootDirectory:    appDir,
		SharedRootDirectory: sharedDir,
		ServiceListeners:    serviceListeners,
		Logger:              logger,
	})
	if err != nil {
		fatal(err)
	}
	accounttypes.RegisterAccountServiceServer(server, service)
	if err := accounttypes.RegisterAccountServiceHandlerServer(ctx, gateway, service); err != nil {
		fatal(err)
	}
	if err := workers.Run(); err != nil && ctx.Err() == nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
