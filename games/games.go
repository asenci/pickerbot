package games

import (
	"errors"
	"strings"

	"github.com/asenci/pickerbot/maps"
)

var (
	GameNotFound      = errors.New("game not found")
	GameAlreadyExists = errors.New("game already exists")
)

type Game struct {
	Maps           maps.Maps
	Name           string
	PlayersPerTeam int
}

func (g *Game) String() string {
	return g.Name
}

type Games map[string]*Game

func (gm Games) Get(name string) (*Game, error) {
	game, ok := gm[strings.ToUpper(name)]

	if !ok {
		return game, GameNotFound
	}

	return game, nil
}

func (gm Games) New(name string, playersPerTeam int, maps maps.Maps) (*Game, error) {
	upperName := strings.ToUpper(name)

	if _, ok := gm[upperName]; ok {
		return nil, GameAlreadyExists
	}

	game := &Game{
		Name:           name,
		PlayersPerTeam: playersPerTeam,
		Maps:           maps,
	}

	gm[strings.ToUpper(name)] = game

	return game, nil
}
