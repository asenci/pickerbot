package cmd

import (
	"fmt"
	"strings"

	"github.com/asenci/pickerbot/games"

	"github.com/Necroforger/dgrouter/exrouter"
)

func Games(ctx *exrouter.Context) {
	var s strings.Builder
	fmt.Fprintln(&s, "Known games:")
	for _, game := range games.All {
		fmt.Fprintf(&s, "  :small_blue_diamond: %s (%d players per team)\n", game.Name, game.PlayersPerTeam)
	}

	ctx.Reply(s.String())
}
