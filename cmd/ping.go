package cmd

import (
	"fmt"

	"github.com/Necroforger/dgrouter/exrouter"
)

// TODO: rate limiting (chill bro)
func Ping(ctx *exrouter.Context) (string, error) {
	return fmt.Sprintf("<@%s> pong", ctx.Msg.Author.ID), nil
}
