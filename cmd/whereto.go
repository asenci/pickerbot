package cmd

import (
	"log"
	"math/rand"
	"strings"

	"github.com/asenci/pickerbot/maps"

	"github.com/asenci/pickerbot/games"

	"github.com/Necroforger/dgrouter/exrouter"
)

func WhereTo(ctx *exrouter.Context) {
	gameName := ctx.Args.Get(1)
	mapName := ctx.Args.Get(2)
	locations := ctx.Args.After(3)

	if locations != "" {
		ls := strings.Split(locations, " ")
		i := rand.Intn(len(ls))
		ctx.Reply("Let's go ", ls[i], "!")
		return
	}

	if gameName == "" {
		ctx.Reply("Which game? Pick one from \"<@", ctx.Ses.State.User.ID, "> games\" or specify a custom one:\n<@", ctx.Ses.State.User.ID, "> whereto <game name> <map name> <location 1> [<location 2> ...]")
		return
	}

	game, err := games.All.Get(gameName)
	if err != nil {
		if err == games.GameNotFound {
			ctx.Reply("I don't know ", gameName, ", give me some locations to draw")
			return
		}

		log.Println(err)
		ctx.Reply("Sorry <@", ctx.Msg.Author.ID, ">, an error has occurred while processing your request: ", err)
	}

	if mapName == "" {
		ctx.Reply("Which map? Pick one from \"<@", ctx.Ses.State.User.ID, "> maps\" or specify a custom one:\n<@", ctx.Ses.State.User.ID, "> whereto <game name> <map name> <location 1> [<location 2> ...]")
		return
	}

	gameMap, err := game.Maps.Get(mapName)
	if err != nil {
		if err == maps.MapNotFound {
			ctx.Reply("I don't know ", mapName, ", give me some locations to draw")
			return
		}

		log.Println(err)
		ctx.Reply("Sorry <@", ctx.Msg.Author.ID, ">, an error has occurred while processing your request: ", err)
	}

	ctx.Reply("Let's go ", gameMap.Locations.Random(), "!")
}
