package main

import (
	"flag"
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
	Token   string
	Prefix  string
	Verbose bool
)

func init() {
	flag.StringVar(&Token, "token", "", "bot token")
	flag.StringVar(&Prefix, "prefix", "", "bot prefix")
	flag.BoolVar(&Verbose, "verbose", false, "increase verbosity")
	flag.Parse()
}

func main() {
	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	router := exrouter.New()
	router.On("draw", cmd.RunDraw).Desc("draw the teams")
	router.On("games", cmd.Games).Desc("list of known games")
	router.On("maps", cmd.Maps).Desc("list of known maps")
	router.On("me", cmd.JoinDraw).Desc("join the latest team draw")
	router.On("ping", cmd.Ping).Desc("responds with pong")
	router.On("play", cmd.NewDraw).Desc("start a team draw")
	router.On("whereto", cmd.WhereTo).Desc("pick a place to drop")
	router.Default = router.On("help", cmd.Help).Desc("prints this help menu")

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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	if Verbose {
		log.Printf("received %s signal, exiting\n", <-sc)
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

	if verbose {
		log.Printf("received request from %s: \"%s\"\n", m.Author.Username, m.Content)
	}

	if err != nil {
		log.Println("error processing request,", err)
	}
}
