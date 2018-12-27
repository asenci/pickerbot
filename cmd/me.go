package cmd

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/asenci/pickerbot/draws"
)

func JoinDraw(ctx *exrouter.Context) {
	draw, err := draws.All.Get(ctx.Msg.ChannelID)
	if err == draws.DrawNotFound {
		ctx.Reply("No draws currently in place, let's start a new one?")
		return
	}
	if err != nil {
		ctx.Reply("Sorry <@", ctx.Msg.Author.ID, ">, an error has occurred while processing your request. Please try again.")
		return
	}

	player := ctx.Msg.Author.ID

	err = draw.Join(player)
	if err == draws.DrawAlreadyJoined {
		ctx.Reply("<@", ctx.Msg.Author.ID, "> you are already on the draw.")
		return
	}
	if err != nil {
		ctx.Reply("Sorry <@", ctx.Msg.Author.ID, ">, an error has occurred while processing your request. Please try again.")
		return
	}

	ctx.Reply("<@", player, "> joined the draw")
}
