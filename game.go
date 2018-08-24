package main

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"sort"
	"strings"
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
	var games []string

	for _, game := range KnownGames {
		games = append(games, game.Name)
	}

	sort.Strings(games)

	msg := &strings.Builder{}
	fmt.Fprintln(msg, "Known games:")
	for _, game := range games {
		fmt.Fprintf(msg, "%s: %d players per team\n", game, KnownGames[game].PlayersPerTeam)
	}

	ctx.Reply(msg)

	return nil
}
