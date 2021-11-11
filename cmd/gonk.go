package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/brandonlbarrow/gonk/internal/discord"
	"github.com/brandonlbarrow/gonk/internal/handler/cocktail"
	"github.com/brandonlbarrow/gonk/internal/handler/info"
	"github.com/sirupsen/logrus"

	"github.com/brandonlbarrow/gonk/internal/handler/stream"

	_ "github.com/joho/godotenv/autoload"
)

var (
	gonkLogLevel      = os.Getenv("GONK_LOG_LEVEL")            // the log level of the main Gonk process. Defaults to Info
	discordgoLogLevel = os.Getenv("DISCORDGO_LOG_LEVEL")       // the log level of the Discordgo session client. See https://pkg.go.dev/github.com/bwmarrin/discordgo#pkg-constants for options. Defaults to LogError
	guildID           = os.Getenv("DISCORD_GUILD_ID")          // the Discord server ID to use for this installation of Gonk.
	channelID         = os.Getenv("DISCORD_STREAM_CHANNEL_ID") // the Discord channel ID to send events for the stream handler to
	token             = os.Getenv("DISCORD_BOT_TOKEN")         // the bot token for use with the Discord API.
	userID            = os.Getenv("DISCORD_USER_ID")           // the Discord user ID to match events on for sending streaming notifications.
	tcdbAPIKey        = os.Getenv("TCDB_API_KEY")              // The Cocktail DB API key for !drank command functionality

	log = newLogger(gonkLogLevel)

	streamHandler = stream.NewHandler(
		stream.WithChannelID(channelID),
		stream.WithGuildID(guildID),
		stream.WithLogger(log),
		stream.WithUserID(userID),
	)
	infoHandler     = info.NewHandler()
	cocktailHandler = cocktail.NewHandler(
		cocktail.WithGuildID(guildID),
		cocktail.WithTCDBAPIKey(tcdbAPIKey),
	)
	handlerMap = map[string]interface{}{
		"stream":   streamHandler.Handle,
		"cocktail": cocktailHandler.Handle,
		"info":     infoHandler.Handle,
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

	log.Debugf("manager created: %v", mgr)

	for name, handler := range handlerMap {
		log.Infof("adding handler %s", name)
		mgr.AddHandler(handler)
	}

	done := make(chan error)

	go run(mgr, done)
	log.Info("GOnk is now running.")
	err := <-done
	if err != nil {
		log.Errorf("GOnk encountered an error while running: %v", err)
		os.Exit(1)
	}
}

func run(mgr *discord.Manager, done chan error) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := mgr.Run(ctx); err != nil {
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
