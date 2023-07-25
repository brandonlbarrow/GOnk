package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/brandonlbarrow/gonk/internal/discord"
	"github.com/brandonlbarrow/gonk/internal/handler"
	"github.com/brandonlbarrow/gonk/internal/handler/cocktail"
	"github.com/brandonlbarrow/gonk/internal/handler/info"
	"github.com/brandonlbarrow/gonk/internal/handler/role"
	"github.com/brandonlbarrow/gonk/v2/internal/webserver"
	"github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload"
)

var (
	gonkLogLevel      = os.Getenv("GONK_LOG_LEVEL")      // the log level of the main Gonk process. Defaults to Info
	discordgoLogLevel = os.Getenv("DISCORDGO_LOG_LEVEL") // the log level of the Discordgo session client. See https://pkg.go.dev/github.com/bwmarrin/discordgo#pkg-constants for options. Defaults to LogError
	token             = os.Getenv("DISCORD_BOT_TOKEN")   // the bot token for use with the Discord API.
	tcdbAPIKey        = os.Getenv("TCDB_API_KEY")        // The Cocktail DB API key for !drank command functionality

	log = newLogger(gonkLogLevel)

	infoHandler     = info.NewHandler()
	cocktailHandler = cocktail.NewHandler(
		cocktail.WithTCDBAPIKey(tcdbAPIKey),
	)
	roleHandler = role.NewHandler()
	handlerMap  = handler.HandlerMap{
		"cocktail": cocktailHandler.Handle,
		"info":     infoHandler.Handle,
		"role":     roleHandler.Handle,
	}
)

func main() {

	// Manager will manage the Discord session, and execute handlers, which will have different behavior based on the server configuration.
	mgr := discord.NewManager(
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go run(ctx, mgr, done)
	go runCallbackServer(ctx, done)
	log.Info("GOnk is now running.")
	err := <-done
	if err != nil {
		log.Errorf("GOnk encountered an error while running: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, mgr *discord.Manager, done chan error) error {
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

func runCallbackServer(ctx context.Context, done chan error) error {
	m := http.NewServeMux()
	webserver.RegisterRoutes(m)
	log.Info("starting webserver")
	if err := http.ListenAndServe(":8080", m); err != nil {
		done <- fmt.Errorf("error running callback server: %w", err)
	}
	return nil
}

func setupTwitch(ctx context.Context) {}
