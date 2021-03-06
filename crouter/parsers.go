package crouter

import (
	"errors"
	"regexp"
	"strings"

	"codeberg.org/eviedelta/drc/detc"
	"github.com/bwmarrin/discordgo"
)

/*
 * This entire file was copy-pasted right from
 * https://codeberg.org/evieDelta/drc/src/branch/master/parse.go
 * licensed under the MIT license
 * https://codeberg.org/evieDelta/drc/src/branch/master/LICENSE.md
**/

var idRegex *regexp.Regexp

// Errors for parsing members
var (
	ErrMemberNotFound  = errors.New("member not found")
	ErrChannelNotFound = errors.New("channel not found")
	ErrRoleNotFound    = errors.New("role not found")
	ErrNoID            = errors.New("input was not an ID")
)

func isID(id string) bool {
	if idRegex == nil {
		idRegex = regexp.MustCompile("([0-9]{15,})")
	}
	return idRegex.MatchString(id)
}

// ParseChannel takes a string and attempts to find a channel that matches that string
func (ctx *Ctx) ParseChannel(channel string) (*discordgo.Channel, error) {
	if isID(channel) {
		// ID Acting mode
		if strings.HasPrefix(channel, "<") {
			if !strings.HasPrefix(channel, "<#") || !strings.HasSuffix(channel, ">") {
				return nil, errors.New("not a channel mention or broken mention")
			}
			channel, _ = detc.Between(channel, "<#", ">")
		}
		return ctx.Session.State.Channel(channel)
	}

	channel = strings.ToLower(channel)

	// Try to find it by name
	g, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil, err
	}

	for _, x := range g.Channels {
		if strings.ToLower(x.Name) == channel {
			return x, nil
		}
	}

	return nil, errors.New("Channel not found")
}

// ParseRole takes a string and attempts to find a role that matches that string
func (ctx *Ctx) ParseRole(role string) (*discordgo.Role, error) {
	if isID(role) {
		// ID Acting mode
		if strings.HasPrefix(role, "<") {
			if !strings.HasPrefix(role, "<@&") || !strings.HasSuffix(role, ">") {
				return nil, errors.New("not a role mention or broken mention")
			}
			role, _ = detc.Between(role, "<@&", ">")
		}
		return ctx.Bot.RoleCache.Get(ctx.Message.GuildID, role)
	}

	role = strings.ToLower(role)

	// Try to find it by name
	g, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil, err
	}

	for _, x := range g.Roles {
		if strings.ToLower(x.Name) == role {
			return x, nil
		}
	}

	return nil, errors.New("Role not found")
}

// ParseMember takes a string and attempts to find a member that matches that string
func (ctx *Ctx) ParseMember(member string) (*discordgo.Member, error) {
	if isID(member) {
		if strings.HasPrefix(member, "<") {
			if !strings.HasPrefix(member, "<@") || !strings.HasSuffix(member, ">") {
				return nil, errors.New("not a member mention or broken mention")
			}
			member, _ = detc.Between(member, "<@", ">")
			if member[0] == '!' {
				member = member[1:]
			}
		}
		m, err := ctx.Bot.MemberCache.Get(ctx.Message.GuildID, member)
		return m, err
	}

	member = strings.ToLower(member)

	// Try to find it by name
	g, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil, err
	}

	for _, x := range g.Members {
		if strings.ToLower(x.User.String()) == member {
			return x, nil
		}
		if strings.ToLower(x.User.Username) == member {
			return x, nil
		}
		if strings.ToLower(x.Nick) == member {
			return x, nil
		}
	}

	return nil, errors.New("Member not found")
}

// ParseUser parses a user
func (ctx *Ctx) ParseUser(user string) (*discordgo.User, error) {
	if m, err := ctx.ParseMember(user); err == nil {
		return m.User, nil
	}

	// try parsing an off-server user
	if !isID(user) {
		return nil, ErrNoID
	}

	return ctx.Session.User(user)
}
