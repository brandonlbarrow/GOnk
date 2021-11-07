package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/brandonlbarrow/gonk/internal/discord"
	"github.com/brandonlbarrow/gonk/internal/handler/cocktail"
	"github.com/sirupsen/logrus"

	"github.com/brandonlbarrow/gonk/internal/handler/stream"

	_ "github.com/joho/godotenv/autoload"
)

var (
	gonkLogLevel      = os.Getenv("GONK_LOG_LEVEL")
	discordgoLogLevel = os.Getenv("DISCORDGO_LOG_LEVEL") // the log level of the Discordgo session client. See https://pkg.go.dev/github.com/bwmarrin/discordgo#pkg-constants for options. Defaults to LogError
	guildID           = os.Getenv("GUILD_ID")            // the Discord server ID to use for this installation of Gonk.
	token             = os.Getenv("DISCORD_BOT_TOKEN")   // the bot token for use with the Discord API.

	log = newLogger(gonkLogLevel)

	streamHandler = stream.NewHandler("", guildID, log)
	handlerMap    = map[string]interface{}{
		"stream":   streamHandler.Handle,
		"cocktail": cocktail.Handler,
	}
)

func main() {

	// TODO remove
	guildID := "308755439145713680"

	mgr := discord.NewManager(
		discord.WithGuildID(guildID),
		discord.MustWithSession(token,
			discord.NewSessionArgs(
				discord.WithLogLevel(
					discord.DiscordLogLevelFromString(discordgoLogLevel)))),
	)

	log.Debug("manager created", mgr)

	for name, handler := range handlerMap {
		logrus.Infof("adding handler %s", name)
		mgr.AddHandler(handler)
	}

	done := make(chan error)

	go run(mgr, done)
	log.Info("GOnk is now running.")
	err := <-done
	if err != nil {
		log.Errorf("GOnk encountered an error while running: %w", err)
		os.Exit(1)
	}
}

func run(mgr *discord.Manager, done chan error) error {

	if err := mgr.Run(context.Background()); err != nil {
		done <- fmt.Errorf("error running manager discordgo session, %w", err)
	}
	return nil
}

func newLogger(level string) *logrus.Logger {

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	l, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(l)
	}
	return logger
}
