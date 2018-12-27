package games

import (
	"fmt"
	"log"
	"testing"

	"github.com/asenci/pickerbot/maps"
)

var (
	emptyGames, populatedGames, largeGames Games
)

func init() {
	emptyGames = Games{}

	populatedGames = Games{
		"GAME1": &Game{Name: "game1"},
	}

	largeGames = Games{}
	for i := 0; i < 1000; i++ {
		_, err := largeGames.New(fmt.Sprintf("game%d", i), 4, maps.Maps{})
		if err != nil {
			log.Panic(err)
		}
	}
}

func BenchmarkGames_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//noinspection GoUnhandledErrorResult
		largeGames.Get("game1")
	}
}

func TestGames_Get(t *testing.T) {
	for _, s := range []string{"game1", "gAme1", "GAME1"} {
		if _, err := populatedGames.Get(s); err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkGames_New(b *testing.B) {
	mm := Games{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//noinspection GoUnhandledErrorResult
		mm.New(fmt.Sprintf("game%d", i), 4, maps.Maps{})
	}
}

func testGamesNew(t *testing.T, mm Games) {
	if _, err := mm.New("new game", 4, maps.Maps{}); err != nil {
		t.Error(err)
	}
}

func TestGames_New(t *testing.T) {
	for _, mm := range []Games{emptyGames, populatedGames} {
		testGamesNew(t, mm)
	}
}
