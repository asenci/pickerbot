package games

import (
	"log"

	"github.com/asenci/pickerbot/maps"
)

var All Games

func initGame(name string, players int, locations map[string][]string) *Game {
	g, err := All.New(name, players, maps.Maps{})
	if err != nil {
		log.Panic(err)
	}

	for m, l := range locations {
		_, err := g.Maps.New(m, l...)
		if err != nil {
			log.Fatal(err)
		}
	}
	return g
}

func init() {
	var g *Game

	All = Games{}

	g = initGame("Apex", 3, map[string][]string{
		"": []string{
			"Airbase",
			"Artillery",
			"Bridges",
			"Bunker",
			"Cascades",
			"Hydro Dam",
			"Market",
			"Relay",
			"Repulsor",
			"Runoff",
			"Skull Town",
			"Slum Lakes",
			"Swamps",
			"The Pit",
			"Thunderdome",
			"Water Treatment",
			"Wetlands",
		},
	})
	g.Name = "Apex Legends"

	g = initGame("BF4", 5, map[string][]string{
		"": []string{"A", "B", "C", "D", "E"},
	})
	g.Name = "Battlefield 4"

	g = initGame("LastTide", 4, map[string][]string{
		"": []string{
			"Ancient Temple",
			"Beattie's Bay",
			"Big Rigs",
			"Eastern Outpost",
			"Eternal Ruins",
			"Exogel Energy",
			"Gellen's Disaster",
			"Gemini Base",
			"Giant's Basin",
			"Good Luck Base",
			"Hades Hideout",
			"Heavy Bomber",
			"Kelp Forest",
			"La Magdalena",
			"Lemone's Triplet",
			"Morg Drilling",
			"New Junk City",
			"Pipetown",
			"Port Reynolds",
			"Red Moon Shipping Co.",
			"Shipping Accident",
			"The Gardens",
			"The Maw",
			"The Spillway",
		},
	})

	g = initGame("Overwatch", 6, map[string][]string{})

	g = initGame("PUBG", 4, map[string][]string{
		"Erangel": []string{
			"Gatka",
			"Georgopol",
			"Kameshki",
			"Lipovka",
			"Mylta",
			"Novorepnoye",
			"Pochinki",
			"Primorsk",
			"Rozhok",
			"Severny",
			"Sosnovka Island",
			"Sosnovka Military Base",
			"Stalber",
			"Yasnaya Polyana",
			"Zharki",
		},
		"Miramar": []string{
			"Chumacera",
			"Cruz del Valle",
			"El Azahar",
			"El Pozo",
			"Hacienda del Patrón",
			"Impala",
			"La Cobreria",
			"Los Higos",
			"Los Leones",
			"Monte Nuevo",
			"Pecado",
			"Prison",
			"Puerto Paraíso",
			"San Martin",
			"Tierra Bronca",
			"Torre Ahumada",
			"Valle del Mar",
		},
		"Sanhok": []string{
			"Ban Tai",
			"Bhan",
			"Camp Alpha",
			"Camp Bravo",
			"Camp Charlie",
			"Cave",
			"Docks",
			"Ha Tinh",
			"Khao",
			"Lakawi",
			"Mongnai",
			"Na Kham",
			"Pai Nan",
			"Quarry",
			"Resort",
			"Ruins",
			"Tambang",
			"Tat Mok",
		},
		"Vikendi": []string{
			"Abbey",
			"Cantra",
			"Castle",
			"Cement Factory",
			"Coal Mine",
			"Cosmodrome",
			"Dino Park",
			"Dobro Mesto",
			"Goroka",
			"Hot Springs",
			"Krichas",
			"Lumber Yard",
			"Milnar",
			"Mount Kreznic",
			"Movatra",
			"Peshkova",
			"Pilnec",
			"Podvosto",
			"Port",
			"Sawmil",
			"Toyar",
			"Trevno",
			"Vihar",
			"Villa",
			"Volnova",
			"Winery",
			"Zabava",
		},
	})
	g.Name = "PlayerUnknown's Battlegrounds"

	g = initGame("RoE", 4, map[string][]string{
		"": []string{
			"Alvitr Castle",
			"Andvari Mine",
			"Balmung City",
			"Bluepeak Town",
			"Cedar Forest",
			"Dione Police Station",
			"Fort Tyrfing",
			"Graveyard",
			"Herschel Academy",
			"Icetongle Village",
			"Lake Herschel",
			"Lumberjack's House",
			"Moose Woods Sawmill",
			"Passer Village",
			"Ring Mountain City",
			"Sidera Lodoicea Ski Resort",
			"Sigel Castle",
			"Skadi City",
			"Snowlake Town",
			"Sober Village",
			"Stele Village",
			"Tiny Village",
			"Valley Town",
			"Wagner City",
			"Whitestone Town",
		},
		"Europa": []string{
			"Angel Town",
			"Artaud City",
			"Bay City",
			"Baywater Construction",
			"Bridge Town",
			"Cassini Farm",
			"Fisherman Town",
			"Galileo Harbor",
			"Marius Observatory",
			"Minos Islands",
			"Ocean Park",
			"Radio Telescope",
			"Resort",
			"Sagan Point",
			"Sarpedon Island",
			"Seafood Market",
			"Seagull Town",
			"Shipyard",
			"Sidon Castle",
			"St. Gabriel Cathedral",
			"Thumb Bay",
			"Tidal Town",
			"Turing City",
			"Wayne Town",
			"Whale Reef",
			"Whiteshore",
			"Windsor Village",
		},
	})
	g.Name = "Ring of Elysium"

}
