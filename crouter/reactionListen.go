package crouter

import "github.com/bwmarrin/discordgo"

// AddReactionHandlerOnce adds a reaction handler function that is only called once
func (ctx *Ctx) AddReactionHandlerOnce(messageID, reaction string, f func(ctx *Ctx)) func() {
	returnFunc := ctx.Session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		if r.UserID != ctx.Message.Author.ID || r.ChannelID != ctx.Channel.ID || r.MessageID != messageID || r.MessageReaction.Emoji.APIName() != reaction {
			return
		}
		f(ctx)

		ctx.Handlers[messageID+reaction]()
		delete(ctx.Handlers, messageID+reaction)

		return
	})
	ctx.Handlers[messageID+reaction] = returnFunc
	return returnFunc
}

// AddReactionHandler adds a reaction handler function
func (ctx *Ctx) AddReactionHandler(messageID, reaction string, f func(ctx *Ctx)) func() {
	returnFunc := ctx.Session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		if r.UserID != ctx.Message.Author.ID || r.ChannelID != ctx.Channel.ID || r.MessageID != messageID || r.MessageReaction.Emoji.APIName() != reaction {
			return
		}
		f(ctx)

		return
	})
	ctx.Handlers[messageID+reaction] = returnFunc
	return returnFunc
}

// AddYesNoHandler reacts with ✅ and ❌, and runs one of two functions depending on which one is used
func (ctx *Ctx) AddYesNoHandler(messageID string, yesFunc, noFunc func(ctx *Ctx)) func() {
	ctx.Session.MessageReactionAdd(ctx.Channel.ID, messageID, "✅")
	ctx.Session.MessageReactionAdd(ctx.Channel.ID, messageID, "❌")

	returnFunc := ctx.Session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		if r.UserID != ctx.Message.Author.ID || r.ChannelID != ctx.Channel.ID || r.MessageID != messageID {
			return
		}

		switch r.MessageReaction.Emoji.APIName() {
		case "✅":
			yesFunc(ctx)
		case "❌":
			noFunc(ctx)
		default:
			return
		}

		ctx.Handlers[messageID]()
		delete(ctx.Handlers, messageID)

		return
	})
	ctx.Handlers[messageID] = returnFunc
	return returnFunc
}
