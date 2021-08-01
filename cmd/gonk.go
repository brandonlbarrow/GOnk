package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/brandonlbarrow/gonk/internal/stream"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func init() {
	switch os.Getenv("ENVIRONMENT") {
	case strings.ToLower("local"):
		if err := godotenv.Load(); err != nil {
			panic(errors.New("No .env file found"))
		}
	}
	return
}

func main() {

	discord := initDiscordSession()
	discord.AddHandler(stream.Handler)

	// https://discord.com/developers/docs/topics/gateway#gateway-intents
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildPresences | discordgo.IntentsGuildMessages)

	err := discord.Open()
	if err != nil {
		fmt.Println("Error opening discord session: ", err)
		return
	}

	fmt.Println("GOnk bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}

func initDiscordSession() *discordgo.Session {

	token, exists := os.LookupEnv("TOKEN")
	if !exists {
		fmt.Println("Cannot find env variable TOKEN. Please ensure this is set to use gonk.")
		os.Exit(1)
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error initializing", err)
		return session
	}

	session.StateEnabled = true

	return session
}
