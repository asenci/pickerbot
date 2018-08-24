package main

import "flag"

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
