package config

import (
	"flag"
	"os"
)

func ParseFlags() *Flags {
	flags := newFlags()

	flag.Var(&flags.Hp, "a", "address and port to run server")
	flag.StringVar(&flags.ShortLink, "b", flags.ShortLink, "address and port to run server")

	flag.Parse()

	if servAdrEnv := os.Getenv("SERVER_ADDRESS"); servAdrEnv != "" {
		flags.ShortLink = servAdrEnv
	}

	if baseURLEnv := os.Getenv("BASE_URL"); baseURLEnv != "" {
		flags.Set(baseURLEnv)
	}

	return &flags
}
