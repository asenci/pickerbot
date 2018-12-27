package maps

import (
	"errors"
	"math/rand"
)

var (
	LocationAlreadyExists = errors.New("location already exists")
)

type Location string

type Locations map[Location]struct{}

func (lm Locations) New(locations ...string) error {
	for _, l := range locations {
		if _, ok := lm[Location(l)]; ok {
			return LocationAlreadyExists
		}

		lm[Location(l)] = struct{}{}
	}

	return nil
}

func (lm Locations) Random() Location {
	i := 0
	r := rand.Intn(len(lm))

	for l := range lm {
		if i == r {
			return l
		}
		i++
	}

	// Keep the compiler happy
	return ""
}
