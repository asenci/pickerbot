package main

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	router := exrouter.New()

	router.On("play", HandlerFuncWrapper(NewDraw)).Desc("start a team draw")
	router.On("games", HandlerFuncWrapper(ListGames)).Desc("list of known games")
	router.On("me", HandlerFuncWrapper(JoinDraw)).Desc("join the latest team draw")
	router.On("ping", HandlerFuncWrapper(Ping)).Desc("responds with pong")

	router.Default = router.On("help", func(ctx *exrouter.Context) {
		helpText := "```"
		for _, v := range router.Routes {
			helpText += fmt.Sprintf("%s: %s\n", v.Name, v.Description)
		}
		helpText += "```"

		ctx.Reply("Hi mate, here is what I can do:", helpText)
	}).Desc("prints this help menu")

	s.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
		if Verbose {
			log.Println("ready to process requests")
		}
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		router.FindAndExecute(s, "", s.State.User.ID, m.Message)
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		router.FindAndExecute(s, "", s.State.User.ID, m.Message)
	})

	if Verbose {
		log.Println("connecting to Discord")
	}
	err = s.Open()
	if err != nil {
		log.Fatal("error connecting to Discord,", err)
	}
	defer s.Close()

	if Verbose {
		defer log.Println("disconnecting from Discord")
	}

	if Verbose {
		log.Println(s.State.User.Username, "connected")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	if Verbose {
		log.Printf("received %s signal, exiting\n", <-sc)
	}
}

func HandlerFuncWrapper(fn func(ctx *exrouter.Context) error) func(ctx *exrouter.Context) {
	return func(ctx *exrouter.Context) {
		if Verbose {
			log.Printf("received request from %s, %s", ctx.Msg.Author.Username, ctx.Args.After(0))
		}

		if err := fn(ctx); err != nil {
			ctx.Reply("Sorry <@", ctx.Msg.Author.ID, ">, an error has occurred while processing your request. Please try again.")
		}
	}
}
