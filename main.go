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

	router.On("play", func(ctx *exrouter.Context) {
		err := NewDraw(ctx)
		if err != nil {
			ctx.Reply("Sorry <@", ctx.Msg.Author.ID, ">, an error has occurred while processing your request. Please try again.")
		}
	}).Desc("start a team draw")

	router.On("me", func(ctx *exrouter.Context) {
		err := JoinDraw(ctx)
		if err != nil {
			ctx.Reply("Sorry <@", ctx.Msg.Author.ID, ">, an error has occurred while processing your request. Please try again.")
		}
	}).Desc(fmt.Sprintf("join the latest team draw"))

	router.On("games", func(ctx *exrouter.Context) {
		ListGames(ctx)
	}).Desc(fmt.Sprintf("list of known games"))

	router.On("ping", func(ctx *exrouter.Context) {
		ctx.Reply("pong")
	}).Desc("responds with pong")

	router.Default = router.On("help", func(ctx *exrouter.Context) {
		helpText := "```"
		for _, v := range router.Routes {
			helpText += fmt.Sprintf("%s: %s\n", v.Name, v.Description)
		}
		helpText += "```"

		ctx.Reply("Hi mate, here is what I can do:", helpText)
	}).Desc("prints this help menu")

	s.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
		log.Println("ready to process requests")
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		router.FindAndExecute(s, "", s.State.User.ID, m.Message)
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		router.FindAndExecute(s, "", s.State.User.ID, m.Message)
	})

	log.Println("connecting to Discord")
	err = s.Open()
	if err != nil {
		log.Fatal("error connecting to Discord,", err)
	}
	defer s.Close()
	defer log.Println("disconnecting from Discord")
	log.Println("connected, press CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	log.Printf("received %s signal, exiting\n", <-sc)
}
