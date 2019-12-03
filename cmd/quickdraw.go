package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/asenci/pickerbot/draws"
	"github.com/asenci/pickerbot/games"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
)

func RunQuickDraw(ctx *exrouter.Context) (string, error) {
	var game *games.Game
	var err error

	// TODO: Move to function so it can be used for both "play" and "quickdraw"
	gameName := ctx.Args.Get(1)
	numPlayersStr := ctx.Args.Get(2)

	if gameName == "" {
		return fmt.Sprintf("Which game? Pick one from **@%s games** or specify a custom one: **@%[1]s quickdraw <__game name__> <__number of players per team__>**", ctx.Ses.State.User.Username), fmt.Errorf("missing game name")
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

	draw := draws.NewDraw(game)

	guild, err := ctx.Ses.Guild(ctx.Msg.GuildID)
	if err != nil {
		return "", err
	}

	presences := guild.Presences

	for _, member := range guild.Members {
		if member.User.Bot {
			continue
		}

		user := member.User

		for _, presence := range presences {
			if presence.User.ID != user.ID {
				continue
			}
			if presence.Status != discordgo.StatusOnline {
				continue
			}
			if err := draw.Join(user.ID); err != nil {
				log.Printf("error adding %q to the draw, %s", user.Username, err)
			}
		}
	}

	teams, err := draw.Run()
	if err != nil {
		if err == draws.DrawNotEnough {
			return "Not enough players for drawing a team", err
		}

		return "", nil
	}

	s := strings.Builder{}
	for _, team := range teams {
		var members []string
		for p := range team.Players {
			members = append(members, p)
		}

		_, err := fmt.Fprintf(&s, "%s: <@%s>\n", team.Name, strings.Join(members, ">, <@"))
		if err != nil {
			return "", err
		}
	}

	return s.String(), nil
}
