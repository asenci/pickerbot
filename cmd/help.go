package cmd

import (
	"fmt"
	"strings"

	"github.com/Necroforger/dgrouter/exrouter"
)

func Help(ctx *exrouter.Context) (string, error) {
	s := strings.Builder{}
	_, err := fmt.Fprintln(&s, "Hi mate, here is what I can do:")
	if err != nil {
		return "", err
	}

	for _, cmd := range ctx.Route.Parent.Routes {
		_, err := fmt.Fprintf(&s, "  :white_small_square: **%s %s**: %s\n", ctx.Ses.State.User.Username, cmd.Name, cmd.Description)
		if err != nil {
			return "", err
		}

	}

	return s.String(), nil
}
