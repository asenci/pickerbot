package main

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	ChannelIdleTimeout = 300 * time.Second
	DrawTimeout        = 15 * time.Second
)

type DrawMap map[string]map[string]*Draw

type Draw struct {
	Game     *Game
	players  []string
	channels []*discordgo.Channel
	timer    *time.Timer
}

func (d *Draw) DrawTeams(ctx *exrouter.Context) {
	game := d.Game.Name
	guild, channel, err := getGuildAndChan(ctx)
	if err != nil {
		ctx.Reply("Error drawing teams for ", game, ", please try again")
		return
	}

	defer delete(currentDraws[guild.ID], channel.ID)

	// Wait for timeout timer to expire
	<-d.timer.C

	players := d.players

	numPlayers := len(players)
	if numPlayers < d.Game.PlayersPerTeam {
		ctx.Reply("Draw failed, not enough players")
		return
	}

	numTeams := (numPlayers / d.Game.PlayersPerTeam) + 1
	ctx.Reply("Drawing ", numTeams, " teams:")

	playersPerTeam := numPlayers / numTeams

	rand.Shuffle(numPlayers, func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	for i := 1; i <= numTeams; i++ {
		teamName := fmt.Sprintf("%s team %d", game, i)

		x := (i - 1) * playersPerTeam
		y := x + playersPerTeam
		teamPlayers := players[x:y]

		ctx.Reply(teamName, ": <@", strings.Join(teamPlayers, ">, <@"), ">")

		teamChan, err := ctx.Ses.GuildChannelCreate(guild.ID, teamName, "voice")
		if err != nil {
			log.Printf("error creating team channel \"%s\", %s", teamName, err)

			ctx.Reply("Error creating channel for team ", teamName, ", please create it manually")
		}
		d.channels = append(d.channels, teamChan)

		for _, player := range teamPlayers {
			err = ctx.Ses.GuildMemberMove(teamChan.GuildID, player, teamChan.ID)
			if err != nil {
				log.Printf("error moving user to channel \"%s\", %s", teamName, err)

				ctx.Reply("Error moving <@", player, "> to channel ", teamName, ", please join manually")
			}
		}

		go func() {
			timer := time.NewTicker(ChannelIdleTimeout)
			for range timer.C {
				if len(teamChan.Recipients) == 0 {
					ctx.Reply("Deleting idle channel ", teamName)

					_, err := ctx.Ses.ChannelDelete(teamChan.ID)
					if err != nil {
						log.Printf("error deleting channel \"%s\", %s", teamName, err)

						ctx.Reply("Error deleting channel ", teamName, ", please delete it manually")
					}
					return
				}
			}
		}()
	}

}

var currentDraws DrawMap

func init() {
	currentDraws = make(DrawMap)
}

func getGuildAndChan(ctx *exrouter.Context) (*discordgo.Guild, *discordgo.Channel, error) {
	channel, err := ctx.Ses.Channel(ctx.Msg.ChannelID)
	if err != nil {
		log.Printf("error getting channel details, %s\n", err)
		return nil, nil, err
	}

	guild, err := ctx.Ses.Guild(channel.GuildID)
	if err != nil {
		log.Printf("error getting guild details, %s\n", err)
		return nil, nil, err
	}

	if _, found := currentDraws[guild.ID]; !found {
		currentDraws[guild.ID] = map[string]*Draw{}
	}

	return guild, channel, nil
}

func NewDraw(ctx *exrouter.Context) error {
	var game *Game

	gameName := ctx.Args.Get(1)
	numPlayersStr := ctx.Args.Get(2)

	if gameName == "" {
		ctx.Reply("Which game? Pick one from \"<@", ctx.Ses.State.User.ID, "> games\" or specify a new one:\n<@", ctx.Ses.State.User.ID, "> play <game name> <number of players per team>")
		return nil
	}

	if numPlayersStr == "" {
		if g, found := KnownGames[strings.ToUpper(gameName)]; !found {
			ctx.Reply("I don't know ", gameName, ", how many players can play on the same team?\nUse: <@", ctx.Ses.State.User.ID, "> play ", gameName, " <number of players per team>")
			return nil
		} else {

			game = g
		}

	} else {
		numPlayers, err := strconv.Atoi(numPlayersStr)
		if err != nil || numPlayers == 0 {
			log.Printf("invalid number: \"%s\", %s\n", numPlayersStr, err)

			ctx.Reply("Invalid number of players: ", numPlayersStr)
			return nil
		}

		game = &Game{gameName, numPlayers}
	}

	guild, channel, err := getGuildAndChan(ctx)
	if err != nil {
		return err
	}

	if draw, found := currentDraws[guild.ID][channel.ID]; found {
		ctx.Reply("Draw for ", draw.Game.Name, " in progress, please wait for it to finish before starting a new one")
		return nil
	}

	ctx.Reply("Sweet! Who is up for some ", gameName, "? (reply with \"<@", ctx.Ses.State.User.ID, "> me\")")

	draw := &Draw{
		Game:  game,
		timer: time.NewTimer(DrawTimeout),
	}
	currentDraws[guild.ID][channel.ID] = draw

	go draw.DrawTeams(ctx)

	return nil
}

func JoinDraw(ctx *exrouter.Context) error {
	guild, channel, err := getGuildAndChan(ctx)
	if err != nil {
		return err
	}

	draw := currentDraws[guild.ID][channel.ID]
	if draw == nil {
		ctx.Reply("No draws currently in place, let's start a new one?")
		return nil
	}

	if !draw.timer.Stop() {
		ctx.Reply("Sorry, too late to join the latest draw. Let's start a new one?")
		return nil
	}

	player := ctx.Msg.Author.ID

	draw.players = append(draw.players, player)
	draw.timer.Reset(DrawTimeout)

	ctx.Reply("<@", player, "> joined the draw")

	return nil
}
