package main

import "github.com/Necroforger/dgrouter/exrouter"

// ToDo: chill bro
func Ping(ctx *exrouter.Context) error {
	ctx.Reply("<@", ctx.Msg.Author.ID, "> pong")
	return nil
}
