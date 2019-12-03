package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Necroforger/dgrouter"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/asenci/pickerbot/cmd"
	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	// ManageChannels bool
	Prefix       string
	Verbose      bool
	PrintVersion bool
	version      = "dev"
)

func init() {
	flag.StringVar(&Token, "token", "", "bot token")
	// flag.BoolVar(&ManageChannels, "manage-channels", false, "manage voice channels")
	flag.StringVar(&Prefix, "prefix", "", "bot prefix")
	flag.BoolVar(&Verbose, "verbose", false, "increase verbosity")
	flag.BoolVar(&PrintVersion, "version", false, "print version")
	flag.Parse()
}

func main() {
	if PrintVersion {
		fmt.Printf("pickerbot %s\n", version)
		return
	}

	if Verbose {
		log.Printf("starting pickerbot %s\n", version)
	}

	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	router := exrouter.New()
	router.On("draw", wrapCmd(cmd.RunDraw, Verbose)).Desc("pick teams for the current draw")
	router.On("games", wrapCmd(cmd.Games, Verbose)).Desc("list known games")
	router.On("maps", wrapCmd(cmd.Maps, Verbose)).Desc("list known game maps")
	router.On("me", wrapCmd(cmd.JoinDraw, Verbose)).Desc("join the current draw")
	router.On("ping", wrapCmd(cmd.Ping, Verbose)).Desc("responds with pong")
	router.On("play", wrapCmd(cmd.NewDraw, Verbose)).Desc("start a new draw")
	router.On("whereto", wrapCmd(cmd.WhereTo, Verbose)).Desc("pick a place to drop or spawn")
	router.Default = router.On("help", wrapCmd(cmd.Help, Verbose)).Desc("prints this help menu")

	s.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
		if Verbose {
			log.Println("ready to process requests")
		}
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		wrapRouter(router, s, m.Message, Prefix, Verbose)
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		wrapRouter(router, s, m.Message, Prefix, Verbose)
	})

	if Verbose {
		log.Println("connecting to Discord")
	}
	err = s.Open()
	if err != nil {
		log.Fatal("error connecting to Discord,", err)
	}
	defer func(s *discordgo.Session) {
		if err := s.Close(); err != nil {
			log.Fatal(err)
		}
	}(s)

	if Verbose {
		defer log.Println("disconnecting from Discord")
	}

	if Verbose {
		log.Println(s.State.User.Username, "connected")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	if Verbose {
		log.Printf("received %s signal, exiting\n", <-sc)
	}
}

func wrapCmd(f func(ctx *exrouter.Context) (string, error), verbose bool) exrouter.HandlerFunc {
	return func(ctx *exrouter.Context) {
		sendReply := func(s string) {
			if verbose {
				log.Printf("sending reply: %q\n", s)
			}
			_, err := ctx.Reply(s)
			if err != nil {
				log.Println("error sending reply,", err)
			}
		}

		reply, err := f(ctx)
		if err != nil {
			log.Println("error processing request,", err)

			if reply == "" {
				reply = fmt.Sprintf("Sorry <@%s>, an error has occurred while processing your request. Please try again.", ctx.Msg.Author.ID)
			}
		}

		if reply != "" {
			sendReply(reply)
		}
	}
}

func wrapRouter(r *exrouter.Route, s *discordgo.Session, m *discordgo.Message, prefix string, verbose bool) {
	// Ignore messages sent by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	err := r.FindAndExecute(s, prefix, s.State.User.ID, m)
	if err == dgrouter.ErrCouldNotFindRoute {
		return
	}

	// TODO: ignore messages without prefix or mention and move logging before execution
	if verbose {
		log.Printf("received request from %s: %q\n", m.Author.Username, m.Content)
	}

	if err != nil {
		log.Println("error processing request,", err)
	}
}
