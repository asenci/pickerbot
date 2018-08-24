package main

import (
	"fmt"

	"github.com/Necroforger/dgrouter/exrouter"
)

func Help(ctx *exrouter.Context) error {
	helpText := "```"
	for _, v := range ctx.Route.Parent.Routes {
		helpText += fmt.Sprintf("%s: %s\n", v.Name, v.Description)
	}
	helpText += "```"

	ctx.Reply("Hi mate, here is what I can do:", helpText)

	return nil
}
