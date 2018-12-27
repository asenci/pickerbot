package cmd

import (
	"fmt"
	"strings"

	"github.com/Necroforger/dgrouter/exrouter"
)

func Help(ctx *exrouter.Context) {
	s := &strings.Builder{}

	fmt.Fprint(s, "```")
	for _, v := range ctx.Route.Parent.Routes {
		fmt.Fprintf(s, "%s: %s\n", v.Name, v.Description)
	}
	fmt.Fprint(s, "```")

	ctx.Reply("Hi mate, here is what I can do:", s)
}
