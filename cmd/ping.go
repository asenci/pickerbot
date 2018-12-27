package cmd

import "github.com/Necroforger/dgrouter/exrouter"

// ToDo: chill bro
func Ping(ctx *exrouter.Context) {
	ctx.Reply("<@", ctx.Msg.Author.ID, "> pong")
}
