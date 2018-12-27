package draws

import (
	"errors"
)

var (
	TeamAlreadyJoined = errors.New("player already joined")
)

type Team struct {
	Name    string
	Players map[string]struct{}
}

func (t *Team) Join(player string) error {
	if _, ok := t.Players[player]; ok {
		return TeamAlreadyJoined
	}

	t.Players[player] = struct{}{}
	return nil
}

type Teams []*Team
