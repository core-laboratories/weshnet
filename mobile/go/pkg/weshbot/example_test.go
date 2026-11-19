package weshbot_test

import (
	"context"
	"os"
	"time"

	qrterminal "github.com/mdp/qrterminal/v3"
	"go.uber.org/zap"
	"moul.io/u"

	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/weshbot"
)

func Example() {
	logger, _ := zap.NewDevelopment()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init bot
	bot, _ := weshbot.New(
		weshbot.WithLogger(logger.Named("botlib")),                                  // configure a logger
		weshbot.WithDisplayName("example bot"),                                      // bot name
		weshbot.WithInsecureMessengerGRPCAddr("127.0.0.1:9091"),                     // connect to running wesh messenger daemon
		weshbot.WithReplay(),                                                        // replay old events as if they are just happening
		weshbot.WithEntityUpdates(),                                                 // include entity updates; i.e., acknowledge
		weshbot.WithFromMyself(),                                                    // trigger hooks for events authored by the bot itself
		weshbot.WithRecipe(weshbot.DebugEventRecipe(logger.Named("debug"))),        // debug events
		weshbot.WithRecipe(weshbot.DelayResponseRecipe(time.Second)),               // add a delay before sending replies
		weshbot.WithRecipe(weshbot.AutoAcceptIncomingContactRequestRecipe()),       // accept incoming contact requests
		weshbot.WithRecipe(weshbot.WelcomeMessageRecipe("welcome to example bot")), // send welcome message to new contacts and new conversations
		weshbot.WithRecipe(weshbot.EchoRecipe("you said: ")),                       // reply to messages with the same message
		weshbot.WithHandler(weshbot.UserMessageHandler, func(ctx weshbot.Context) { // custom handler
			ctx.ReplyString("hello world!")
		}),
	)

	// display link and qr code
	logger.Info("retrieve instance Wesh ID", zap.String("pk", bot.PublicKey()), zap.String("link", bot.WeshIDURL()))
	qrterminal.GenerateHalfBlock(bot.WeshIDURL(), qrterminal.L, os.Stdout)

	// signal handling
	go func() {
		u.WaitForCtrlC()
		cancel()
	}()

	// start bot
	bot.Start(ctx)
}
