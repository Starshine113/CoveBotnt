package modcommands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/starshine-sys/covebotnt/cbdb"
	"github.com/starshine-sys/covebotnt/crouter"
	"github.com/starshine-sys/covebotnt/etc"
	"github.com/starshine-sys/covebotnt/structs"
	"github.com/starshine-sys/flagparser"
)

// LogMute adds a mute to the mute log
func LogMute(ctx *crouter.Ctx) (err error) {
	if len(ctx.Args) < 1 {
		ctx.CommandError(&crouter.ErrorNotEnoughArgs{
			NumRequiredArgs: 1,
			SuppliedArgs:    len(ctx.Args),
		})
		return nil
	}

	flagParser, err := flagparser.NewFlagParser(flagparser.Duration("d", "dur", "duration"), flagparser.Bool("hardmute"))
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	args, err := flagParser.Parse(ctx.Args)
	if err != nil {
		ctx.CommandError(err)
		return nil
	}
	if len(args["rest"].([]string)) == 0 {
		ctx.CommandError(&crouter.ErrorMissingRequiredArgs{
			RequiredArgs: "<user ID/mention>",
			MissingArgs:  "<user ID/mention>",
		})
		return nil
	}
	remaining := args["rest"].([]string)

	member, err := ctx.ParseMember(remaining[0])
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	reason := "None"
	if len(remaining) > 1 {
		reason = strings.Join(remaining[1:], " ")
	}
	muteType := "mute"
	if args["hardmute"].(bool) {
		muteType = "hardmute"
	}
	defaultDuration, _ := time.ParseDuration("876600h")
	var duration string
	if args["d"].(time.Duration) == defaultDuration {
		duration = "none"
	} else {
		duration = etc.HumanizeDuration(etc.DurationPrecisionSeconds, args["d"].(time.Duration))
	}

	entry, err := ctx.Database.AddToModLog(&cbdb.ModLogEntry{
		GuildID:   ctx.Message.GuildID,
		UserID:    member.User.ID,
		ModID:     ctx.Author.ID,
		Type:      muteType,
		Reason:    reason + fmt.Sprintf(" (duration: %v)", duration),
		Time:      time.Now().UTC(),
		Snowflake: ctx.Bot.SnowflakeGen.Get(),
	})
	if err != nil {
		ctx.CommandError(err)
		return nil
	}

	_, err = ctx.Send(fmt.Sprintf("%v Added this mute to the mod log.", crouter.SuccessEmoji))
	if err != nil {
		return err
	}

	logEmbed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    ctx.Author.String(),
			IconURL: ctx.Author.AvatarURL("256"),
		},
		Color:       0xe5da00,
		Title:       fmt.Sprintf("User %vd (case #%v)", muteType, entry.ID),
		Description: fmt.Sprintf("**%v** (ID: %v)", member.User.String(), member.User.ID),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: member.User.AvatarURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Reason",
				Value:  reason,
				Inline: false,
			},
			{
				Name:   "Duration",
				Value:  duration,
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Moderator ID: " + ctx.Author.ID,
		},
		Timestamp: entry.Time.Format(time.RFC3339),
	}

	modLog := ctx.AdditionalParams["guildSettings"].(*structs.GuildSettings).Moderation.ModLog

	if modLog == "" {
		_, err = ctx.Send(fmt.Sprintf("%v No mod log channel set. Set one with `%vmodlog <channel>`.", crouter.WarnEmoji, ctx.GuildPrefix))
		return
	}

	_, err = ctx.Session.ChannelMessageSendEmbed(modLog, logEmbed)
	return
}
