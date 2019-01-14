package cmd

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/asenci/pickerbot/maps"

	"github.com/asenci/pickerbot/games"

	"github.com/Necroforger/dgrouter/exrouter"
)

func WhereTo(ctx *exrouter.Context) (string, error) {
	gameName := ctx.Args.Get(1)
	mapName := ctx.Args.Get(2)
	locations := ctx.Args.After(3)

	if locations != "" {
		ls := strings.Split(locations, " ")
		i := rand.Intn(len(ls))
		return fmt.Sprintf("Let's go %q!", ls[i]), nil
	}

	if gameName == "" {
		return fmt.Sprintf("Which game? Pick one from **@%[1]s games** or specify a custom one: **@%[1]s whereto <__game name__> <__map name__> <__location 1__> <__location 2__> ...**", ctx.Ses.State.User.Username), fmt.Errorf("missing game name")
	}

	game, err := games.All.Get(gameName)
	if err != nil {
		if err == games.GameNotFound {
			return fmt.Sprintf("I don't know the game %q, give me some locations to draw: **@%s whereto %[1]s %[3]s <__location 1__> <__location 2__> ...**", gameName, ctx.Ses.State.User.Username, mapName), err
		}

		return "", err
	}

	gameMap, err := game.Maps.Get(mapName)
	if err != nil {
		if err == maps.MapNotFound {
			if mapName == "" {
				return fmt.Sprintf("Which map? Pick one from **@%s maps %s** or specify a custom one: **@%[1]s whereto %[2]s <__map name__> <__location 1__> <__location 2__> ...**", ctx.Ses.State.User.Username, gameName), fmt.Errorf("missing map name")
			}

			return fmt.Sprintf("I don't know the map %q, give me some locations to draw: **@%s whereto %s %[1]s <__location 1__> <__location 2__> ...**", mapName, ctx.Ses.State.User.Username, gameName), err
		}

		return "", err
	}

	return fmt.Sprintf("Let's go %q!", gameMap.Locations.Random()), nil
}
