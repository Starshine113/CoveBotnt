package main

import (
	"github.com/bwmarrin/discordgo"
)

const successEmoji, errorEmoji string = "✅", "❌"

func commandTree(command string, args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error

	switch command {
	case "ping":
		err = commandPing(s, m)
	case "help":
		err = commandHelp(args, s, m)
	case "setstatus":
		err = commandSetStatus(args, s, m)
	case "starboard":
		err = commandStarboard(args, s, m)
	case "modroles":
		err = commandModRoles(args, s, m)
	case "echo":
		err = commandEcho(args, s, m)
	case "prefix":
		err = commandPrefix(args, s, m)
	}

	if err != nil {
		sugar.Errorf("Error running command %v", err)
	}
}
