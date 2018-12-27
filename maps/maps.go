package maps

import (
	"errors"
	"strings"
)

var (
	MapAlreadyExists = errors.New("map already exists")
	MapNotFound      = errors.New("map not found")
)

type Map struct {
	Name      string
	Locations Locations
}

func (m *Map) String() string {
	return m.Name
}

type Maps map[string]*Map

func (mm Maps) Get(name string) (*Map, error) {
	m, ok := mm[strings.ToUpper(name)]

	if !ok {
		m, ok = mm[""]
		if !ok {
			return m, MapNotFound
		}
	}

	return m, nil
}

func (mm Maps) New(name string, locations ...string) (*Map, error) {
	if _, ok := mm[strings.ToUpper(name)]; ok {
		return nil, MapAlreadyExists
	}

	m := &Map{name, Locations{}}
	if err := m.Locations.New(locations...); err != nil {
		return nil, err
	}
	mm[strings.ToUpper(name)] = m

	return m, nil
}
