package cmd

import (
	"fmt"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/asenci/pickerbot/draws"
)

func JoinDraw(ctx *exrouter.Context) (string, error) {
	draw, err := draws.All.Get(ctx.Msg.ChannelID)
	if err != nil {
		if err == draws.DrawNotFound {
			return fmt.Sprintf("No draws currently in place, let's start a new one? Pick a game from **@%s games**", ctx.Ses.State.User.Username), err
		}

		return "", err
	}

	err = draw.Join(ctx.Msg.Author.ID)
	if err != nil {
		if err == draws.DrawAlreadyJoined {
			return fmt.Sprintf("<@%s> you have already joined the draw.", ctx.Msg.Author.ID), err
		}

		return "", err
	}

	return fmt.Sprintf("<@%s> joined the draw", ctx.Msg.Author.ID), nil
}
