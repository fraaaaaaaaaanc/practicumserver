package config

import (
	"errors"
	"strconv"
	"strings"
)

// Hp describes the server address data structure.
type Hp struct {
	Host string // Host is the server's hostname or IP address.
	Port int    // Port is the port on which the server listens.
}

// Flags defines the fields to store flag values used in the application's configuration.
type Flags struct {
	Prefix          string // Prefix is the base URL for the application.
	LogLevel        string // LogLevel is the logging level.
	FileLog         bool   // FileLog enables or disables file logging.
	ConsoleLog      bool   // ConsoleLog enables or disables console logging (not used in the code).
	FileStoragePath string // FileStoragePath is the path to the file storage.
	DBStorageAdress string // DBStorageAdress is the database connection string.
	Hp                     // Embed the Hp struct to handle server address configuration.
}

// newFlags initializes the Flags struct with default values.
func newFlags() Flags {
	return Flags{
		Prefix: "http://localhost:8080",
		Hp: Hp{
			Host: "localhost",
			Port: 8080,
		},
		LogLevel:        "info",
		FileLog:         false,
		FileStoragePath: "",
		DBStorageAdress: "",
	}
}

// String method overrides the default String method to format the server address as a string.
func (h *Hp) String() string {
	return h.Host + ":" + strconv.Itoa(h.Port)
}

// Set method overrides the default Set method to validate and set the server address.
func (h *Hp) Set(addres string) error {
	if addres == "" {
		h.Host = "localhost"
		h.Port = 8080
		return nil
	}
	hp := strings.Split(addres, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	h.Host = hp[0]
	h.Port = port
	return nil
}
