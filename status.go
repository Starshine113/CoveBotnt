package main

import (
	"github.com/bwmarrin/discordgo"
)

func onReady(s *discordgo.Session, _ *discordgo.Ready) {
	sugar.Name = s.State.User.Username
	sugar.AvatarURL = s.State.User.AvatarURL("128")

	err := updateStatus(s)
	if err != nil {
		sugar.Errorf("Error updating status: %v", err)
	}
}

func updateStatus(s *discordgo.Session) (err error) {
	switch config.Bot.CustomStatus.Type {
	case "listening":
		err = s.UpdateListeningStatus(config.Bot.CustomStatus.Status)
	case "playing":
		err = s.UpdateGameStatus(0, config.Bot.CustomStatus.Status)
	case "watching":
		err = s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{{
				Name: config.Bot.CustomStatus.Status,
				Type: 3,
			}},
		})
	default:
		return nil
	}
	return err
}
