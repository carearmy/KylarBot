package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/bwmarrin/discordgo"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

var token string

func main() {
	dummyConfig()

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

	// Register userJoin as a callback for the guildMemberAdd events.
	dg.AddHandler(userJoin)

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

func userJoin(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	guild, _ := s.Guild(event.GuildID)
	s.GuildMemberRoleAdd(event.GuildID, event.Member.User.ID, getRole(guild, "Civilian").ID)
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

	if m.Content == "!check_roles" {
		member, err := s.State.Member(g.ID, m.Author.ID)
		if err != nil {
			return
		}

		var buffer bytes.Buffer

		for i := 0; i < len(member.Roles); i++ {
			mRole, err := s.State.Role(g.ID, member.Roles[i])
			if err != nil {
				return
			}

			json_m, _ := json.Marshal(mRole)
			buffer.WriteString(string(json_m) + "\n")
		}

		s.ChannelMessageSend(c.ID, buffer.String())
	}

	if strings.HasPrefix(m.Content, "!math ") {
		var expression *govaluate.EvaluableExpression
		var err error
		var result interface{}

		if expression, err = govaluate.NewEvaluableExpression(m.Content[len("!math "):]); err != nil {
			s.ChannelMessageSend(c.ID, fmt.Sprintf("<@!%s> Something went wrong when calculating the results: %s", m.Author.ID, err.Error()))
			return
		}
		if result, err = expression.Evaluate(nil); err != nil {
			s.ChannelMessageSend(c.ID, fmt.Sprintf("<@!%s> Something went wrong when calculating the results: %s", m.Author.ID, err.Error()))
			return
		}

		s.ChannelMessageSend(c.ID, fmt.Sprintf("<@!%s> Result: %v", m.Author.ID, result))
	}

	if m.Content == "!ping" {
		s.ChannelMessageSend(c.ID, "Pong.")
	}

	if m.Content == "!pong" {
		s.ChannelMessageSend(c.ID, "Ping.")
	}

	if m.Content == "!req_chan" {
		allowed, err := hasRole(s, g, m.Author, "Member")
		if err != nil {
			fmt.Println("Something went wrong when checking the roles for an author: ", err)
			return
		}

		if allowed {
			s.ChannelMessageSend(c.ID, fmt.Sprintf("<@!%s> Not implemented.", m.Author.ID))
		} else {
			s.ChannelMessageSend(c.ID, fmt.Sprintf("<@!%s> You don't have permission to use this command.", m.Author.ID))
		}
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

	/* Mute while in-dev.
	for _, channel := range event.Guild.Channels {
		_, _ = s.ChannelMessageSend(channel.ID, "Hello everyone. My name's Kylar, and soon I'll be able to do a lot of awesome stuff!")
	}*/
}

func getRole(g *discordgo.Guild, role string) *discordgo.Role {
	for _, r := range g.Roles {
		if r.Name == role {
			return r
		}
	}

	return nil
}

func hasRole(s *discordgo.Session, g *discordgo.Guild, user *discordgo.User, role string) (bool, error) {

	var member *discordgo.Member
	var checkRole *discordgo.Role
	var err error

	if member, err = s.State.Member(g.ID, user.ID); err != nil {
		return false, err
	}

	if checkRole = getRole(g, role); checkRole == nil {
		return false, errors.New("Role not found")
	}

	for i := 0; i < len(member.Roles); i++ {
		var mRole *discordgo.Role
		if mRole, err = s.State.Role(g.ID, member.Roles[i]); err != nil {
			return false, err
		}

		if mRole.Name == role || checkRole.Position < mRole.Position {
			return true, nil
		}
	}

	return false, nil

}
