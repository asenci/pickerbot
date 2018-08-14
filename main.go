package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
)

type MatchRequest struct {
	channels       []*discordgo.Channel
	game           string
	players        []string
	playersPerTeam int
	timer          *time.Timer
}

var (
	Token  string
	Prefix string

	matchRequests = map[string]*MatchRequest{}
)

func init() {
	flag.StringVar(&Token, "token", "", "bot token")
	flag.StringVar(&Prefix, "prefix", "!", "bot prefix")
	flag.Parse()
}

func main() {
	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	s.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
		log.Println("ready to process requests")
	})

	router := exrouter.New()

	router.On("ping", func(ctx *exrouter.Context) {
		fmt.Printf("%v+\n", ctx)
		ctx.Reply("pong")
	}).Desc("responds with pong")

	router.On("play", func(ctx *exrouter.Context) {
		var (
			game           string
			playersStr     string
			playersPerTeam int
			err            error
		)

		game = ctx.Args.Get(1)
		playersStr = ctx.Args.Get(2)

		if game == "" {
			ctx.Reply("Which game?")
			return
		}

		if playersStr == "" {
			if v, ok := MaxPlayersPerTeam[game]; ok {
				playersPerTeam = v
			} else {
				ctx.Reply("I don't know ", game, ", how many players can play on the same team?")
				return
			}
		} else {
			playersPerTeam, err = strconv.Atoi(playersStr)
			if err != nil || playersPerTeam == 0 {
				ctx.Reply("Invalid number of players: ", playersStr)
			}
		}

		matchRequests[ctx.Msg.ChannelID] = &MatchRequest{
			game:           game,
			playersPerTeam: playersPerTeam,
			timer:          time.NewTimer(5 * time.Second),
		}
		ctx.Reply("Sweet! Who is up for some ", game, "? (Reply with \"@TestBot me\")")
		router.On("me", func(ctx *exrouter.Context) {
			matchRequests[ctx.Msg.ChannelID].timer.Reset(5 * time.Second)
			matchRequests[ctx.Msg.ChannelID].players = append(matchRequests[ctx.Msg.ChannelID].players, ctx.Msg.Author.ID)
		}).Desc(fmt.Sprintf("enter the team draw for %s", game))

		go func(context *exrouter.Context) {
			mr := matchRequests[context.Msg.ChannelID]
			<-mr.timer.C

			players := mr.players
			numPlayers := len(players)
			playersPerTeam := mr.playersPerTeam

			numTeams := numPlayers / playersPerTeam
			if numTeams == 0 {
				numTeams++
			}

			if numTeams == 1 {
				ctx.Reply("Drawing ", numTeams, " team.")
			} else {
				ctx.Reply("Drawing ", numTeams, " teams.")
			}

			rand.Shuffle(numPlayers, func(i, j int) {
				players[i], players[j] = players[j], players[i]
			})
			for i := 1; i <= numTeams; i++ {
				teamName := fmt.Sprintf("%s team %d", mr.game, i)

				x := (i - 1) * playersPerTeam
				y := x + playersPerTeam
				if y > numPlayers {
					y = numPlayers
				}
				teamPlayers := players[x:y]

				ctx.Reply(fmt.Sprintf("%s: <@%s>", teamName, strings.Join(teamPlayers, ">, <@")))

				channel, err := ctx.Ses.Channel(ctx.Msg.ChannelID)
				if err != nil {
					log.Printf("error getting channel details, %s\n", err)
					continue
				}

				teamChan, err := ctx.Ses.GuildChannelCreate(channel.GuildID, teamName, "voice")
				if err != nil {
					log.Printf("error creating team channel \"%s\", %s", teamName, err)
				}
				mr.channels = append(mr.channels, teamChan)

				for _, player := range teamPlayers {
					err = ctx.Ses.GuildMemberMove(teamChan.GuildID, player, teamChan.ID)
					if err != nil {
						log.Printf("error moving user to channel \"%s\", %s", teamName, err)
					}
				}

				go func(c *discordgo.Channel) {
					timer := time.NewTicker(300 * time.Second)
					for range timer.C {
						if len(c.Recipients) == 0 {
							_, err := ctx.Ses.ChannelDelete(teamChan.ID)
							if err != nil {
								log.Printf("error deleting channel \"%s\", %s", teamName, err)
							}
							timer.Stop()
						}
					}
				}(teamChan)
			}
		}(ctx)
	}).Desc("draw players for a match")

	router.Default = router.On("help", func(ctx *exrouter.Context) {
		helpText := ""
		helloText := fmt.Sprintln("Hi mate, here is what I can do:")

		helpText += "```"
		for _, v := range router.Routes {
			helpText += fmt.Sprintf("%s: %s\n", v.Name, v.Description)
		}
		helpText += "```"

		ctx.Reply(helloText, helpText)
	}).Desc("prints this help menu")

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		router.FindAndExecute(s, "", s.State.User.ID, m.Message)
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		router.FindAndExecute(s, "", s.State.User.ID, m.Message)
	})

	log.Println("connecting to Discord")
	err = s.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	defer s.Close()
	defer log.Println("disconnecting from Discord")
	log.Println("connected, press CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	log.Printf("received %s signal, exiting\n", <-sc)
}
