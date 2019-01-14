package cmd

import (
	"log"
	"strconv"
	"time"

	"github.com/asenci/pickerbot/draws"
	"github.com/asenci/pickerbot/games"

	"github.com/Necroforger/dgrouter/exrouter"
)

func NewDraw(ctx *exrouter.Context) {
	var game *games.Game
	var err error

	gameName := ctx.Args.Get(1)
	numPlayersStr := ctx.Args.Get(2)

	if gameName == "" {
		ctx.Reply("Which game? Pick one from **@", ctx.Ses.State.User.Username, " games** or specify a custom one: **@", ctx.Ses.State.User.Username, " play <__game name__> <__number of players per team__>**")
		return
	}

	if numPlayersStr == "" {
		game, err = games.All.Get(gameName)
		if err != nil {
			if err == games.GameNotFound {
				ctx.Reply("I don't know *", gameName, "*, how many players can play on the same team? Use: **@", ctx.Ses.State.User.Username, " play ", gameName, " <__number of players per team__>**")
				return
			}
		}

	} else {
		numPlayers, err := strconv.Atoi(numPlayersStr)
		if err != nil || numPlayers == 0 {
			log.Printf("invalid number: \"%s\", %s\n", numPlayersStr, err)

			ctx.Reply("Invalid number of players: ", numPlayersStr)
			return
		}

		game = &games.Game{
			Name:           gameName,
			PlayersPerTeam: numPlayers,
			Maps:           nil,
		}
	}

	_, err = draws.All.New(ctx.Msg.ChannelID, game)
	if err == draws.DrawAlreadyExists {
		ctx.Reply("Draw in progress, please wait for it to finish before starting a new one")
		return
	}
	if err != nil {
		ctx.Reply("Sorry <@", ctx.Msg.Author.ID, ">, an error has occurred while processing your request. Please try again.")
		return
	}

	ctx.Reply("Sweet! Who is up for some ", game, "? Reply with **@", ctx.Ses.State.User.Username, " me**")

	time.AfterFunc(5*time.Minute, func() {
		_, err := draws.All.Get(ctx.Msg.ChannelID)
		if err == draws.DrawNotFound {
			return
		}

		RunDraw(ctx)
	})
}
