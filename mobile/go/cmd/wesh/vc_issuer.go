package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/oklog/run"
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/weshvcissuer"
	"github.com/core-laboratories/weshnet/v2/pkg/verifiablecredstypes"
)

func vcIssuerCommand() *ffcli.Command {
	listenerFlag := ""
	serverRootURL := ""
	allowedRedirectURIs := ""
	vcIssuerPrivKey := ""

	fsBuilder := func() (*flag.FlagSet, error) {
		fs := flag.NewFlagSet("wesh vc-issuer", flag.ExitOnError)
		fs.StringVar(&listenerFlag, "http.listener", listenerFlag, "http listener")
		fs.StringVar(&serverRootURL, "http.server-root", serverRootURL, "http server root")
		fs.StringVar(&allowedRedirectURIs, "http.allowed-redirect-uris", allowedRedirectURIs, "comma-separated allowed callback URLs; tone://vc is always allowed")
		fs.StringVar(&vcIssuerPrivKey, "vc-sk", vcIssuerPrivKey, "Verifiable Credentials Issuer private key (base64 encoded)")
		return fs, nil
	}

	return &ffcli.Command{
		Name:           "vc-issuer",
		ShortUsage:     "wesh [global flags] vc-issuer [flags]",
		ShortHelp:      "start a verified credentials issuer service",
		Options:        ffSubcommandOptions(),
		FlagSetBuilder: fsBuilder,
		UsageFunc:      usageFunc,
		Exec: func(ctx context.Context, args []string) error { //nolint:revive
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			if listenerFlag == "" {
				return fmt.Errorf("missing -http.listener flag")
			}

			if serverRootURL == "" {
				return fmt.Errorf("missing -http.server-root flag")
			}

			if vcIssuerPrivKey == "" {
				return fmt.Errorf("missing -vc-sk flag")
			}

			g := run.Group{}
			g.Add(func() error {
				<-ctx.Done()
				return ctx.Err()
			}, func(error) {
				cancel()
			})

			logger, err := manager.GetLogger()
			if err != nil {
				return err
			}

			skBytes, err := base64.RawStdEncoding.DecodeString(vcIssuerPrivKey)
			if err != nil {
				return err
			}

			if len(skBytes) != ed25519.SeedSize {
				return fmt.Errorf("invalid sk size")
			}

			issuerPrivateKey := &[32]byte{}
			copy(issuerPrivateKey[:], skBytes)

			l, err := (&net.ListenConfig{}).Listen(ctx, "tcp", listenerFlag)
			if err != nil {
				return err
			}

			vcConfig := &weshvcissuer.Config{
				ServerRootURL: serverRootURL,
				IssuerSignKey: issuerPrivateKey,
				Flow: &weshvcissuer.FlowConfig{
					Type:          verifiablecredstypes.FlowType_FlowTypeCode, // TODO
					CodeGenerator: weshvcissuer.CodeGeneratorZero,            // TODO
					CodeSenderClient: &weshvcissuer.PhoneCodeSenderMockService{
						Logger: logger,
					}, // TODO
				},
				Logger: logger,
			}
			if allowedRedirectURIs != "" {
				vcConfig.AllowedRedirectURIs = strings.Split(allowedRedirectURIs, ",")
			}

			issuer, err := weshvcissuer.New(vcConfig)
			if err != nil {
				return err
			}

			server := &http.Server{
				Handler:           issuer,
				ReadHeaderTimeout: time.Second * 5,
			}

			g.Add(func() error {
				return server.Serve(l)
			}, func(_ error) {
				l.Close()
			})

			return g.Run()
		},
	}
}
