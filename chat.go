package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Chatter - container for all things Gonk can say
type Chatter struct {
	Responses map[string]string
	Commands  []string
	About     map[string]string
	Queue     chan string
}

func chatHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	chatter := bootstrap()
	fmt.Println(m.Content)

	if m.Content == "help" {
		chatter.Queue <- m.Content
		chatter.gonkChatHandler(s, m.Message)
	}

	if m.Content == "whoami" {
		user, err := s.User(m.Author.ID)
		if err != nil {
			fmt.Println(err)
		}
		s.ChannelMessageSend(m.ChannelID, user.Mention())
	}

	// for _, usrMention := range m.Mentions {
	// 	if usrMention.Username == "gonk" {
	// 		chatter.gonkChatHandler(s, m.Message)
	// 	}
	// }
}

func (c Chatter) gonkChatHandler(s *discordgo.Session, m *discordgo.Message) {

	message := <-c.Queue
	fmt.Println("gonkChatHandler says " + message)

	switch message {
	case "":
		s.ChannelMessageSend(m.ChannelID, c.About["gif"])
		s.ChannelMessageSend(m.ChannelID, "~I am GOnk. Type help for functions.~")
	case "help":
		s.ChannelMessageSend(m.ChannelID, "My commands are "+strings.Join(c.Commands, "\n"))
	}
}

func bootstrap() Chatter {

	responses := make(map[string]string)
	responses["back me up here @gonk"] = "~he does have a point~"
	about := make(map[string]string)
	about["gif"] = "https://www.micechat.com/wp-content/uploads/2019/05/EG-series-power-droid-gif.gif"
	queue := make(chan string)
	chatter := Chatter{
		responses,
		[]string{"whoami", "multi", "magic8ball"},
		about,
		queue,
	}

	return chatter
}

func listMembers(s *discordgo.Session, g string) []*discordgo.Member {

	memberList, err := s.GuildMembers(g, "0", 100)
	if err != nil {
		fmt.Println(err)
	}

	return memberList
}
