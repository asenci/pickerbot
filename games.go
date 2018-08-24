package main

import (
	"fmt"
	"strings"

	"github.com/Necroforger/dgrouter/exrouter"
)

type Game struct {
	Name           string
	PlayersPerTeam int
}

func (g *Game) String() string {
	return g.Name
}

type GameRegistry map[string]*Game

func (m GameRegistry) Register(name string, playerPerTeam int) {
	m[strings.ToUpper(name)] = &Game{
		Name:           name,
		PlayersPerTeam: playerPerTeam,
	}
}

var KnownGames GameRegistry

func init() {
	KnownGames = GameRegistry{}
	KnownGames.Register("BF4", 5)
	KnownGames.Register("Overwatch", 6)
	KnownGames.Register("PUBG", 4)
}

func ListGames(ctx *exrouter.Context) error {
	var msg strings.Builder
	fmt.Fprintln(&msg, "Known games:")
	for _, game := range KnownGames {
		fmt.Fprintf(&msg, "%s: %d players per team\n", game.Name, game.PlayersPerTeam)
	}

	ctx.Reply(msg.String())

	return nil
}
