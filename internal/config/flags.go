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

	if logLvlEnv := os.Getenv("LOG_LEVEL"); logLvlEnv != "" {
		flags.LogLevel = logLvlEnv
	}
	if fileStorageFile := os.Getenv("FILE_STORAGE_PATH"); fileStorageFile != "" {
		flags.FileStorage = fileStorageFile
	}
}

func ParseFlags() *Flags {
	flags := newFlags()

	flag.Var(&flags.Hp, "a", "address and port to run server")
	flag.StringVar(&flags.ShortLink, "b", flags.ShortLink, "address and port to run server")
	flag.StringVar(&flags.LogLevel, "l", flags.LogLevel, "log level")
	flag.BoolVar(&flags.FileLog, "fl", flags.FileLog, "On file logging")
	flag.StringVar(&flags.FileStorage, "f", flags.FileStorage, "On file storage")

	flag.Parse()

	return &flags
}
