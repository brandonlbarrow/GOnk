package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/brandonlbarrow/gonk/internal/discord"
	"github.com/brandonlbarrow/gonk/internal/handler/cocktail"
	"github.com/sirupsen/logrus"

	"github.com/brandonlbarrow/gonk/internal/handler/stream"
)

var (
	streamHandler = &stream.Handler{}
	handlerMap    = map[string]interface{}{
		"stream":   streamHandler.Handle,
		"cocktail": cocktail.Handler,
	}

	discordgoLogLevel = os.Getenv("DISCORDGO_LOG_LEVEL") // the log level of the Discordgo session client. See https://pkg.go.dev/github.com/bwmarrin/discordgo#pkg-constants for options. Defaults to LogError
	guildID           = os.Getenv("GUILD_ID")            // the Discord server ID to use for this installation of Gonk.
	token             = os.Getenv("DISCORD_BOT_TOKEN")   // the bot token for use with the Discord API.
)

func main() {

	token, exists := os.LookupEnv("TOKEN")
	if !exists {
		fmt.Println("Cannot find env variable TOKEN. Please ensure this is set to use gonk.")
		os.Exit(1)
	}

	// TODO remove
	guildID := "308755439145713680"

	mgr := discord.NewManager(
		discord.WithGuildID(guildID),
		discord.MustWithSession(token, discord.NewSessionArgsWithDefaults()),
	)

	for name, handler := range handlerMap {
		logrus.Infof("adding handler %s", name)
		mgr.AddHandler(handler)
	}

	if err := mgr.Run(context.Background()); err != nil {
		fmt.Errorf("error %w", err)
		return
	}

	fmt.Println("GOnk bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

func run() error {

}
