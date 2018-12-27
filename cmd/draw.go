package cmd

import (
	"log"
	"strings"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/asenci/pickerbot/draws"
)

func RunDraw(ctx *exrouter.Context) {
	defer func(channel string) {
		err := draws.All.Delete(channel)
		if err != nil {
			log.Println("error removing draw,", err)
		}
	}(ctx.Msg.ChannelID)

	draw, err := draws.All.Get(ctx.Msg.ChannelID)
	if err == draws.DrawNotFound {
		ctx.Reply("No draws currently in place, let's start a new one?")
		return
	}
	if err != nil {
		ctx.Reply("An error has occurred while processing your request. Please try again.")
		return
	}

	teams, err := draw.Run()
	if err == draws.DrawNotEnough {
		ctx.Reply("Not enough players for drawing a team")
		return
	}
	if err != nil {
		ctx.Reply("An error has occurred while processing your request. Please try again.")
		return
	}

	channel, err := ctx.Ses.Channel(ctx.Msg.ChannelID)
	if err != nil {
		log.Printf("failed to retrieve channel, %s", err)
	}

	guild := channel.GuildID

	for _, team := range teams {
		var members []string
		for p, _ := range team.Players {
			members = append(members, p)
		}

		ctx.Reply(team.Name, ": <@", strings.Join(members, ">, <@"), ">")

		teamChan, err := ctx.Ses.GuildChannelCreate(guild, team.Name, "voice")
		if err != nil {
			log.Printf("error creating team channel %q, %s", team.Name, err)

			ctx.Reply("Error creating channel for team ", team.Name, ", please create it manually")
		}

		for _, player := range members {
			err = ctx.Ses.GuildMemberMove(guild, player, teamChan.ID)
			if err != nil {
				log.Printf("error moving user to channel %q, %s", team.Name, err)

				ctx.Reply("Error moving <@", player, "> to channel ", team.Name, ", please join manually")
			}
		}

		// TODO: cleanup idle channels (use Guild.VoiceStates?)
		// go func(channel, name string) {
		// 	timer := time.NewTicker(5 * time.Minute)
		// 	for range timer.C {
		// 		c, err := ctx.Ses.Channel(channel)
		// 		if err != nil {
		// 			log.Printf("failed to find channel %q", name)
		//
		// 			ctx.Reply("Error deleting channel ", name, ", please delete it manually")
		// 		}
		//
		// 		if len(c.Recipients) == 0 {
		// 			ctx.Reply("Deleting idle channel ", name)
		//
		// 			_, err := ctx.Ses.ChannelDelete(channel)
		// 			if err != nil {
		// 				log.Printf("error deleting channel %q, %s", name, err)
		//
		// 				ctx.Reply("Error deleting channel ", name, ", please delete it manually")
		// 			}
		// 			return
		// 		}
		// 	}
		// }(teamChan.ID, team.Name)
	}
}
