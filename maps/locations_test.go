package maps

import (
	"fmt"
	"testing"
)

var emptyLocations, populatedLocations, largeLocations Locations

func init() {
	emptyLocations = Locations{}

	populatedLocations = Locations{
		"location1": struct{}{},
	}

	largeLocations = Locations{}
	for i := 0; i < 1000; i++ {
		largeLocations[Location(fmt.Sprintf("location%d", i))] = struct{}{}
	}
}

func BenchmarkLocations_New(b *testing.B) {
	lm := Locations{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//noinspection GoUnhandledErrorResult
		lm.New(fmt.Sprintf("location%d", i))
	}
}

func testLocationsNew(t *testing.T, lm Locations) {
	if err := lm.New("new location"); err != nil {
		t.Error(err)
	}
}

func TestLocations_New(t *testing.T) {
	for _, l := range []Locations{emptyLocations, populatedLocations} {
		testLocationsNew(t, l)
	}
}

func BenchmarkLocations_Random(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		largeLocations.Random()
	}
}

func testLocationsRandom(t *testing.T, lm Locations) {
	t.Logf("Random location: %+v", lm.Random())
}

func TestLocations_Random(t *testing.T) {
	for _, l := range []Locations{emptyLocations, populatedLocations} {
		testLocationsRandom(t, l)
	}
}
