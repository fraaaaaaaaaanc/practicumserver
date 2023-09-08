package config

import "flag"

type Flags struct {
	HostFlag  string
	ShortLink string
}

func ParseFlags() *Flags {
	flags := &Flags{}

	flag.StringVar(&flags.HostFlag, "a", ":8080", "address and port to run server")
	flag.StringVar(&flags.ShortLink, "b", "http://localhost:8080/", "address and port to run server")

	flag.Parse()

	return flags
}
