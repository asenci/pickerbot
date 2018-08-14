package main

import "flag"

var (
	Token  string
	Prefix string
)

func init() {
	flag.StringVar(&Token, "token", "", "bot token")
	flag.StringVar(&Prefix, "prefix", "", "bot prefix")
	flag.Parse()
}
