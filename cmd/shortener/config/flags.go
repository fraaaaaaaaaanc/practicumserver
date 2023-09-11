package config

import (
	"flag"
)

func ParseFlags() *Flags {
	flags := newFlags()

	flag.Var(&flags.Hp, "a", "address and port to run server")
	flag.StringVar(&flags.ShortLink, "b", flags.ShortLink, "address and port to run server")

	flag.Parse()

	return &flags
}
