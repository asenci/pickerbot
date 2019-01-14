package cmd

import (
	"fmt"
	"strings"

	"github.com/asenci/pickerbot/games"

	"github.com/Necroforger/dgrouter/exrouter"
)

func Games(ctx *exrouter.Context) (string, error) {
	s := strings.Builder{}
	_, err := fmt.Fprintln(&s, "Known games:")
	if err != nil {
		return "", err
	}

	for alias, game := range games.All {
		_, err := fmt.Fprintf(&s, "  :video_game: %s (%d players per team): **@%s play %s**\n", game.Name, game.PlayersPerTeam, ctx.Ses.State.User.Username, strings.ToLower(alias))
		if err != nil {
			return "", err
		}
	}

	return s.String(), nil
}
