package draws

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/asenci/pickerbot/games"
)

var (
	DrawAlreadyExists = errors.New("draw already exists")
	DrawAlreadyJoined = errors.New("player already joined")
	DrawNotEnough     = errors.New("not enough players")
	DrawNotFound      = errors.New("draw not found")
)

type Draw struct {
	game    *games.Game
	players map[string]struct{}
}

func (d *Draw) Run() (Teams, error) {
	numPlayers := len(d.players)
	if numPlayers <= d.game.PlayersPerTeam {
		return nil, DrawNotEnough
	}

	players := make([]string, 0, numPlayers)
	for p := range d.players {
		players = append(players, p)
	}
	rand.Shuffle(numPlayers, func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	numTeams := numPlayers / d.game.PlayersPerTeam
	if numPlayers%d.game.PlayersPerTeam > 0 {
		numTeams++
	}

	teams := make(Teams, 0, numTeams)

	playersPerTeam := numPlayers / numTeams
	leftovers := numPlayers % numTeams

	var first, last int
	for i := 0; i < numTeams; i++ {
		team := &Team{
			fmt.Sprintf("Team %d", i+1),
			make(map[string]struct{}),
		}

		first = last
		last = first + playersPerTeam
		if i < leftovers {
			last++
		}

		for _, player := range players[first:last] {
			if err := team.Join(player); err != nil {
				return nil, err
			}
		}

		teams = append(teams, team)
	}

	return teams, nil
}

func (d *Draw) Join(player string) error {
	if _, ok := d.players[player]; ok {
		return DrawAlreadyJoined
	}

	d.players[player] = struct{}{}
	return nil
}

type Draws map[string]*Draw

func (dm Draws) Delete(channel string) {
	delete(dm, channel)
}

func (dm Draws) Get(channel string) (*Draw, error) {
	draw, ok := dm[channel]
	if !ok {
		return nil, DrawNotFound
	}

	return draw, nil
}

func (dm Draws) New(channel string, game *games.Game) (*Draw, error) {

	if _, ok := dm[channel]; ok {
		return nil, DrawAlreadyExists
	}

	d := &Draw{
		game:    game,
		players: map[string]struct{}{},
	}
	dm[channel] = d

	return d, nil
}
