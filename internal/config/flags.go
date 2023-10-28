// Package config provides a configuration mechanism for the application.
// The package defines the Flags struct, which holds various configuration parameters for the application.
// It also includes functions for parsing configuration flags from the command line and reading configuration
// values from environment variables. These functions allow users to configure the application's behavior.
package config

import (
	"flag"
	"os"
)

// ParseConfFlags parses configuration flags. It first checks the flags passed during program execution.
// Then, it calls ParseEnv to check environment variables and updates the values in the Flags struct accordingly.
// Finally, it returns the Flags struct with the values set from flags and environment variables.
func ParseConfFlags() *Flags {
	flags := ParseFlags()
	ParseEnv(flags)
	return flags
}

// ParseEnv checks the presence of values in environment variables. If the environment variables are not empty,
// their values are assigned to the fields of the Flags struct, potentially overriding values set by the user during
// program execution.
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

// ParseFlags defines and parses command-line flags.
func ParseFlags() *Flags {
	flags := newFlags()

	flag.Var(&flags.Hp, "a", "address and port to run server")
	flag.StringVar(&flags.Prefix, "b", flags.Prefix, "prefix for response")
	flag.StringVar(&flags.LogLevel, "l", flags.LogLevel, "log level")
	flag.BoolVar(&flags.FileLog, "fl", flags.FileLog, "On file logging")
	flag.StringVar(&flags.FileStoragePath, "f", flags.FileStoragePath, "On file storage")
	flag.StringVar(&flags.DBStorageAdress, "d", flags.DBStorageAdress, "On db params: host, user, password, dbname")

	flag.Parse()

	return &flags
}
