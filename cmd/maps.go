package cmd

import (
	"fmt"
	"strings"

	"github.com/asenci/pickerbot/games"

	"github.com/Necroforger/dgrouter/exrouter"
)

func Maps(ctx *exrouter.Context) (string, error) {
	gameName := ctx.Args.Get(1)

	if gameName == "" {
		return fmt.Sprintf("Which game? Pick one from **@%s games**", ctx.Ses.State.User.Username), fmt.Errorf("missing game name")
	}

	game, err := games.All.Get(gameName)
	if err != nil {
		if err == games.GameNotFound {
			return fmt.Sprintf("I don't know the game %q. Pick one from **@%s games**", gameName, ctx.Ses.State.User.Username), err
		}

		return "", err
	}

	s := strings.Builder{}
	_, err = fmt.Fprintf(&s, "Known %q maps:\n", game)
	if err != nil {
		return "", err
	}

	for _, m := range game.Maps {
		var f string
		if m.Name == "" {
			f = "  :earth_americas: %s[Default map]\n"
		} else {
			f = "  :earth_americas: %s\n"
		}

		_, err := fmt.Fprintf(&s, f, m)
		if err != nil {
			return "", err
		}
	}

	return s.String(), nil
}
