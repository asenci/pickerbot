package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/asenci/pickerbot/draws"
	"github.com/asenci/pickerbot/games"

	"github.com/Necroforger/dgrouter/exrouter"
)

func NewDraw(ctx *exrouter.Context) (string, error) {
	var game *games.Game
	var err error

	gameName := ctx.Args.Get(1)
	numPlayersStr := ctx.Args.Get(2)

	if gameName == "" {
		return fmt.Sprintf("Which game? Pick one from **@%s games** or specify a custom one: **@%[1]s play <__game name__> <__number of players per team__>**", ctx.Ses.State.User.Username), fmt.Errorf("missing game name")
	}

	if numPlayersStr == "" {
		game, err = games.All.Get(gameName)
		if err != nil {
			if err == games.GameNotFound {
				return fmt.Sprintf("I don't know the game %q, how many players can play on the same team? **@%s play %[1]s <__number of players per team__>**", gameName, ctx.Ses.State.User.Username), err
			}

			return "", err
		}

	} else {
		numPlayers, err := strconv.Atoi(numPlayersStr)
		if err != nil || numPlayers < 2 {
			return fmt.Sprintf("Invalid number of players: %s", numPlayersStr), fmt.Errorf("invalid number %q, %s", numPlayersStr, err)
		}

		game = &games.Game{
			Name:           gameName,
			PlayersPerTeam: numPlayers,
			Maps:           nil,
		}
	}

	_, err = draws.All.New(ctx.Msg.ChannelID, game)
	if err != nil {
		if err == draws.DrawAlreadyExists {
			return "Draw in progress, please wait for it to finish before starting a new one", err
		}

		return "", err
	}

	time.AfterFunc(5*time.Minute, func() {
		_, err := draws.All.Get(ctx.Msg.ChannelID)
		if err != nil {
			if err == draws.DrawNotFound {
				return
			}

			log.Println("error retrieving the draw,", err)
		}

		draws.All.Delete(ctx.Msg.ChannelID)
		_, err = ctx.Reply("Draw timed out. Please start a new one")
		if err != nil {
			log.Println("error sending reply,", err)
		}
	})

	return fmt.Sprintf("Sweet! Who is up for some %q? Use **@%s me** to join the draw and **@%[2]s draw** to pick the teams when ready", game, ctx.Ses.State.User.Username), nil
}
