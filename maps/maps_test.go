package maps

import (
	"fmt"
	"log"
	"testing"
)

var (
	emptyMaps, populatedMaps, largeMaps Maps
)

func init() {
	emptyMaps = Maps{}

	populatedMaps = Maps{
		"MAP1": &Map{Name: "map1"},
	}

	largeMaps = Maps{}
	for i := 0; i < 1000; i++ {
		m, err := largeMaps.New(fmt.Sprintf("map%d", i))
		if err != nil {
			log.Panic(err)
		}

		for j := 0; j < 1000; j++ {
			err := m.Locations.New(fmt.Sprintf("location%d", i))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func BenchmarkMaps_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//noinspection GoUnhandledErrorResult
		largeMaps.Get("map1")
	}
}

func TestMaps_Get(t *testing.T) {
	for _, s := range []string{"map1", "mAp1", "MAP1"} {
		if _, err := populatedMaps.Get(s); err != nil {
			t.Error(err)
		}
	}
}

func TestMaps_Get_Default(t *testing.T) {
	mm := Maps{}
	_, err := mm.New("", "location1")
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range []string{"map1", ""} {
		if _, err := mm.Get(s); err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkMaps_New(b *testing.B) {
	mm := Maps{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//noinspection GoUnhandledErrorResult
		mm.New(fmt.Sprintf("map%d", i))
	}
}

func testMapsNew(t *testing.T, mm Maps) {
	if _, err := mm.New("new map"); err != nil {
		t.Error(err)
	}
}

func TestMaps_New(t *testing.T) {
	for _, mm := range []Maps{emptyMaps, populatedMaps} {
		testMapsNew(t, mm)
	}
}
