package main

import (
	"log"
	"strings"

	"github.com/mb-14/gomarkov"

	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s.State.User.ID == m.Author.ID {
		return
	}

	content := m.ContentWithMentionsReplaced()

	if strings.HasPrefix(content, "-markov") {
		// Give them a chain!
		tokens := []string{gomarkov.StartToken}
		for tokens[(len(tokens)-1)] != gomarkov.EndToken {
			next, _ := chain.Generate(tokens[(len(tokens) - 1):])
			tokens = append(tokens, next)
		}

		s.ChannelMessageSend(m.ChannelID, strings.Join(tokens[1:len(tokens)-1], " "))
		return
	}

	split := strings.Split(content, " ")

	chain.Add(split)
}

func discordReady(s *discordgo.Session, _ *discordgo.Ready) {
	log.Printf("[info] ready! logged in as %s\n", s.State.User.String())
	saveModel()
}
