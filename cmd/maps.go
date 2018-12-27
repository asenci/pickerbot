package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/asenci/pickerbot/games"

	"github.com/Necroforger/dgrouter/exrouter"
)

func Maps(ctx *exrouter.Context) {
	gameName := ctx.Args.Get(1)

	if gameName == "" {
		ctx.Reply("Which game? Pick one from \"<@", ctx.Ses.State.User.ID, "> games\"")
		return
	}

	game, err := games.All.Get(gameName)
	if err != nil {
		if err == games.GameNotFound {
			ctx.Reply("Sorry <@", ctx.Msg.Author.ID, ">, I don't know ", gameName)
			return
		}

		log.Println(err)
		ctx.Reply("Sorry <@", ctx.Msg.Author.ID, ">, an error has occurred while processing your request: ", err)
	}

	var s strings.Builder
	fmt.Fprintf(&s, "Known %q maps:\n", game)
	for _, m := range game.Maps {
		if m.Name == "" {
			fmt.Fprintf(&s, "  :small_blue_diamond: %s[Default map]\n", m)
		} else {
			fmt.Fprintf(&s, "  :small_blue_diamond: %s\n", m)
		}
	}

	ctx.Reply(s.String())
}
