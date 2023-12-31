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
		flags.Prefix = baseURLEnv
	}

	if servAdrEnv := os.Getenv("SERVER_ADDRESS"); servAdrEnv != "" {
		flags.Set(servAdrEnv)
	}

	if LogLvlEnv := os.Getenv("LOG_LEVEL"); LogLvlEnv != "" {
		flags.LogLevel = LogLvlEnv
	}
	if fileStorageFileEnv := os.Getenv("FILE_STORAGE_PATH"); fileStorageFileEnv != "" {
		flags.FileStorage = fileStorageFileEnv
	}
	if dbAdressEnv := os.Getenv("DATABASE_DSN"); dbAdressEnv != "" {
		flags.DBAdress = dbAdressEnv
	}
}

func ParseFlags() *Flags {
	flags := newFlags()

	flag.Var(&flags.Hp, "a", "address and port to run server")
	flag.StringVar(&flags.Prefix, "b", flags.Prefix, "address and port to run server")
	flag.StringVar(&flags.LogLevel, "l", flags.LogLevel, "log level")
	flag.BoolVar(&flags.FileLog, "fl", flags.FileLog, "On file logging")
	flag.StringVar(&flags.FileStorage, "f", flags.FileStorage, "On file storage")
	flag.StringVar(&flags.DBAdress, "d", flags.DBAdress, "On file storage")

	flag.Parse()

	return &flags
}
