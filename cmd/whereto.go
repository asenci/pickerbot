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
		ctx.Reply("Let's go *", ls[i], "*!")
		return
	}

	if gameName == "" {
		ctx.Reply("Which game? Pick one from **@", ctx.Ses.State.User.Username, " games** or specify a custom one: **@", ctx.Ses.State.User.Username, " whereto <__game name__> <__map name__> <__location 1__> <__location 2__> ...**")
		return
	}

	game, err := games.All.Get(gameName)
	if err != nil {
		if err == games.GameNotFound {
			ctx.Reply("I don't know *", gameName, "*, give me some locations to draw")
			return
		}

		log.Println(err)
		ctx.Reply("Sorry <@", ctx.Msg.Author.ID, ">, an error has occurred while processing your request: ", err)
	}

	gameMap, err := game.Maps.Get(mapName)
	if err != nil {
		if err == maps.MapNotFound {
			if mapName == "" {
				ctx.Reply("Which map? Pick one from **@", ctx.Ses.State.User.Username, " maps ", gameName, "** or specify a custom one: **@", ctx.Ses.State.User.Username, " whereto <__game name__> <__map name__> <__location 1__> <__location 2__> ...**")
				return
			}

			ctx.Reply("I don't know *", mapName, "*, give me some locations to draw")
			return
		}

		log.Println(err)
		ctx.Reply("Sorry <@", ctx.Msg.Author.ID, ">, an error has occurred while processing your request: ", err)
	}

	ctx.Reply("Let's go *", gameMap.Locations.Random(), "*!")
}
