package main

import (
	"math/rand"
	"strings"

	"github.com/Necroforger/dgrouter/exrouter"
)

func WhereTo(ctx *exrouter.Context) error {
	gameName := ctx.Args.Get(1)
	mapName := ctx.Args.Get(2)
	gameNameUpper := strings.ToUpper(gameName)
	mapNameUpper := strings.ToUpper(mapName)
	options := strings.Split(ctx.Args.After(2), " ")

	if len(options) == 0 {
		if gameNameUpper == "PUBG" {
			if mapNameUpper == "ERANGEL" {
				options = []string{"Gatka", "Georgopol", "Kameshki", "Lipovka", "Mylta", "Novorepnoye", "Pochinki", "Primorsk", "Rozhok", "Severny", "Sosnovka Island", "Sosnovka Military Base", "Stalber", "Stalber", "Yasnaya Polyana", "Zharki"}
			} else if mapNameUpper == "MIRAMAR" {
				options = []string{"Chumacera", "Cruz del Valle", "El Azahar", "El Pozo", "Hacienda del Patrón", "Impala", "La Cobreria", "Los Higos", "Los Leones", "Monte Nuevo", "Pecado", "Prison", "Puerto Paraíso", "San Martin", "Tierra Bronca", "Torre Ahumada", "Valle del Mar"}
			} else if mapNameUpper == "SANHOK" {
				options = []string{"Ban Tai", "Bhan", "Camp Alpha", "Camp Alpha", "Camp Bravo", "Camp Bravo", "Camp Charlie", "Camp Charlie", "Cave", "Docks", "Ha Tinh", "Khao", "Lakawi", "Mongnai", "Na Kham", "Pai Nan", "Quarry", "Resort", "Ruins", "Tambang", "Tat Mok"}
			} else {
				ctx.Reply("I don't know ", mapName, "give me some options")
			}
		} else if gameNameUpper == "BF4" {
			options = []string{"A", "B", "C", "D", "E"}
		} else {
			ctx.Reply("I don't know ", gameName, "give me some options")
		}
	}

	i := rand.Intn(len(options))
	ctx.Reply("Let's go ", options[i], "!")

	return nil
}
