package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	qrterminal "github.com/mdp/qrterminal/v3"
	"github.com/oklog/run"
	"go.uber.org/zap"
	"moul.io/u"
	"moul.io/zapconfig"

	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/weshbot"
	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/weshversion"
)

var (
	username     = u.CurrentUsername("anon")
	node1Addr    = flag.String("addr1", "127.0.0.1:9091", "first remote 'wesh daemon' address")
	node2Addr    = flag.String("addr2", "127.0.0.1:9092", "second remote 'wesh daemon' address")
	displayName1 = flag.String("name1", username+" (testbot1)", "first bot's display name")
	displayName2 = flag.String("name2", username+" (testbot2)", "second bot's display name")
	debug        = flag.Bool("debug", false, "debug mode")
	logFormat    = flag.String("log-format", "console", strings.Join(zapconfig.AvailablePresets, ", "))
)

func main() {
	flag.Parse()
	if err := Main(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		os.Exit(1)
	}
}

type TestBot struct {
	Bot1, Bot2 *weshbot.Bot
	ctx        context.Context
	logger     *zap.Logger
}

func Main() error {
	config := zapconfig.Configurator{}
	config.SetPreset(*logFormat)
	logger := config.MustBuild()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init testbot
	testbot := &TestBot{ctx: ctx, logger: logger}
	if err := testbot.InitBot1(); err != nil {
		return fmt.Errorf("init bot 1: %w", err)
	}
	if err := testbot.InitBot2(); err != nil {
		return fmt.Errorf("init bot 2: %w", err)
	}

	// FIXME: bot1 and bot2 should be contact

	// start testbot
	var g run.Group
	g.Add(func() error { return testbot.Bot1.Start(ctx) }, func(error) { cancel() })
	g.Add(func() error { return testbot.Bot2.Start(ctx) }, func(error) { cancel() })
	g.Add(run.SignalHandler(ctx, syscall.SIGKILL))
	return g.Run()
}

func (testbot *TestBot) VersionCommand(ctx weshbot.Context) {
	_ = ctx.ReplyString("version: " + weshversion.Version)
	// FIXME: also returns the version of the remote messenger and protocol
}

// InitBot1 initializes the entrypoint bot
func (testbot *TestBot) InitBot1() error {
	logger := testbot.logger.Named("bot1")

	// init bot
	opts := []weshbot.NewOption{}
	opts = append(opts,
		weshbot.WithLogger(logger.Named("lib")),                                  // configure a logger
		weshbot.WithDisplayName(*displayName1),                                   // bot name
		weshbot.WithInsecureMessengerGRPCAddr(*node1Addr),                        // connect to running wesh messenger daemon
		weshbot.WithRecipe(weshbot.AutoAcceptIncomingContactRequestRecipe()),    // accept incoming contact requests
		weshbot.WithRecipe(weshbot.WelcomeMessageRecipe("welcome to testbot1")), // send welcome message to new contacts and new conversations
		weshbot.WithRecipe(weshbot.EchoRecipe("you said1: ")),                   // reply to messages with the same message
		// FIXME: with auto-send `/help` suggestion on welcome
		weshbot.WithCommand("version", "show version", testbot.VersionCommand),
	)
	if *debug {
		opts = append(opts, weshbot.WithRecipe(weshbot.DebugEventRecipe(logger.Named("debug")))) // debug events
	}
	bot, err := weshbot.New(opts...)
	if err != nil {
		return fmt.Errorf("bot initialization failed: %w", err)
	}
	// display link and qr code
	logger.Info("retrieve instance Wesh ID",
		zap.String("pk", bot.PublicKey()),
		zap.String("link", bot.WeshIDURL()),
	)
	qrterminal.GenerateHalfBlock(bot.WeshIDURL(), qrterminal.L, os.Stdout)
	testbot.Bot1 = bot
	return nil
}

// InitBot2 initializes the companion bot
func (testbot *TestBot) InitBot2() error {
	logger := testbot.logger.Named("bot2")

	// init bot
	opts := []weshbot.NewOption{}
	opts = append(opts,
		weshbot.WithLogger(logger.Named("lib")),                                  // configure a logger
		weshbot.WithDisplayName(*displayName2),                                   // bot name
		weshbot.WithInsecureMessengerGRPCAddr(*node2Addr),                        // connect to running wesh messenger daemon
		weshbot.WithRecipe(weshbot.AutoAcceptIncomingContactRequestRecipe()),    // accept incoming contact requests
		weshbot.WithRecipe(weshbot.WelcomeMessageRecipe("welcome to testbot2")), // send welcome message to new contacts and new conversations
		weshbot.WithRecipe(weshbot.EchoRecipe("you said2: ")),                   // reply to messages with the same message
	)
	if *debug {
		opts = append(opts, weshbot.WithRecipe(weshbot.DebugEventRecipe(logger.Named("debug")))) // debug events
	}
	bot, err := weshbot.New(opts...)
	if err != nil {
		return fmt.Errorf("bot initialization failed: %w", err)
	}
	// display link and qr code
	logger.Info("retrieve instance Wesh ID",
		zap.String("pk", bot.PublicKey()),
		zap.String("link", bot.WeshIDURL()),
	)
	qrterminal.GenerateHalfBlock(bot.WeshIDURL(), qrterminal.L, os.Stdout)
	testbot.Bot2 = bot
	return nil
}
