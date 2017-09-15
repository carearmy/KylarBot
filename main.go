package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

var token string

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events
	dg.AddHandler(ready)

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Register guildCreate as a callback for the guildCreate events.
	dg.AddHandler(guildCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	fmt.Println("KylarBot is now running.  Press CTRL-C to exit.")
	// Simple way to keep the program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Make sure we're not idle
	_ = s.UpdateStatus(0, "")
}

// This function will be called every time a new message is created on any channel the bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Find the channel that the message came from.
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return
	}

	if m.Content == "!k" {
		s.ChannelMessageSend(c.ID, "That's me!")
	}

	if m.Content == "!ping" {
		s.ChannelMessageSend(c.ID, "Pong.")
	}

	if m.Content == "!pong" {
		s.ChannelMessageSend(c.ID, "Ping.")
	}

	if strings.HasPrefix(m.Content, "!setgame") {
		allowed, err := hasRole(s, g, m.Author, "Administrator")
		if err != nil {
			fmt.Println("Something went wrong when checking the roles for an author: ", err)
			return
		}

		if allowed {
			s.UpdateStatus(0, m.Content[len("!setgame"):])
		} else {
			s.ChannelMessageSend(c.ID, fmt.Sprintf("<@!%s> You don't have permission to use this command.", m.Author.ID))
		}
	}
}

// This function will be called every time a newguild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable == true {
		return
	}

	for _, channel := range event.Guild.Channels {
		_, _ = s.ChannelMessageSend(channel.ID, "Hello everyone. My name's Kylar, and soon I'll be able to do a lot of awesome stuff!")
	}
}

func hasRole(s *discordgo.Session, g *discordgo.Guild, user *discordgo.User, role string) (bool, error) {
	member, err := s.State.Member(g.ID, user.ID)
	if err != nil {
		return false, err
	}

	for i := 0; i < len(member.Roles); i++ {
		mRole, err := s.State.Role(g.ID, member.Roles[i])
		if err != nil {
			return false, err
		}

		return mRole.Name == role, nil
	}

	return false, nil
}
