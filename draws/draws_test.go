package draws

import (
	"fmt"
	"testing"

	"github.com/asenci/pickerbot/games"
)

var (
	largeDraw Draw
)

func init() {
	largeDraw = Draw{
		game: &games.Game{
			Name:           "game",
			PlayersPerTeam: 4,
		},
		players: map[string]struct{}{},
	}
	for i := 0; i < 1000; i++ {
		largeDraw.players[fmt.Sprintf("player %d", i)] = struct{}{}
	}
}

func BenchmarkDraw_Run(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := largeDraw.Run()
		if err != nil {
			b.Error(err)
		}
	}
}

func TestDraw_Run(t *testing.T) {
	for i := 1; i <= 8; i++ {
		for j := i + 1; j <= 16; j++ {
			draw := Draw{
				game: &games.Game{
					Name:           "game",
					PlayersPerTeam: i,
				},
				players: map[string]struct{}{},
			}

			for k := 1; k <= j; k++ {
				if err := draw.Join(fmt.Sprintf("player %d", k)); err != nil {
					t.Error(err)
				}
			}

			teams, err := draw.Run()
			if err != nil {
				t.Errorf("error drawing %d-player teams for %d players, %s", i, j, err)
			}

			for player := range draw.players {
				var found bool

				for _, team := range teams {
					for p := range team.Players {
						if _, ok := draw.players[p]; !ok {
							t.Errorf("unknown player %q", player)
						}

						if p == player {
							found = true
						}
					}
				}

				if !found {
					t.Errorf("player %q not found in any of the teams", player)
				}
			}
		}
	}
}
