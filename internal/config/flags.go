package config

import (
	"flag"
	"os"
)

func ParseConfFlugs() *Flags {
	flags := ParseFlags()
	ParseEnv(flags)
	return flags
}

func ParseEnv(flags *Flags) {
	if baseURLEnv := os.Getenv("BASE_URL"); baseURLEnv != "" {
		flags.ShortLink = baseURLEnv
	}

	if servAdrEnv := os.Getenv("SERVER_ADDRESS"); servAdrEnv != "" {
		flags.Set(servAdrEnv)
	}
}

func ParseFlags() *Flags {
	flags := newFlags()

	flag.Var(&flags.Hp, "a", "address and port to run server")
	flag.StringVar(&flags.ShortLink, "b", flags.ShortLink, "address and port to run server")

	flag.Parse()

	return &flags
}