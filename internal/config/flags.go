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
		flags.FileStoragePath = fileStorageFileEnv
	}
	if dbAdressEnv := os.Getenv("DATABASE_DSN"); dbAdressEnv != "" {
		flags.DBStorageAdress = dbAdressEnv
	}
}

func ParseFlags() *Flags {
	flags := newFlags()

	flag.Var(&flags.Hp, "a", "address and port to run server")
	flag.StringVar(&flags.Prefix, "b", flags.Prefix, "address and port to run server")
	flag.StringVar(&flags.LogLevel, "l", flags.LogLevel, "log level")
	flag.BoolVar(&flags.FileLog, "fl", flags.FileLog, "On file logging")
	flag.StringVar(&flags.FileStoragePath, "f", flags.FileStoragePath, "On file storage")
	flag.StringVar(&flags.DBStorageAdress, "d", flags.DBStorageAdress, "On db params: host, user, password, dbname")

	flag.Parse()

	return &flags
}
