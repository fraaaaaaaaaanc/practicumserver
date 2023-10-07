package config

import (
	"errors"
	"strconv"
	"strings"
)

type Hp struct {
	Host string
	Port int
}
type Flags struct {
	Hp
	Prefix      string
	LogLevel    string
	FileLog     bool
	ConsoleLog  bool
	FileStorage string
	DBAdress    string
}

func newFlags() Flags {
	return Flags{
		Prefix: "http://localhost:8080",
		Hp: Hp{
			Host: "localhost",
			Port: 8080,
		},
		LogLevel:    "info",
		FileLog:     false,
		FileStorage: "C:/Users/frant/go/go1.21.0/bin/pkg/mod/github.com/fraaaaaaaaaanc/practicumserver/internal/tmp/short-url-db.json",
		//FileStorage: "/tmp/short-url-db.json",
		DBAdress: "host=localhost user=postgres password=1234 dbname=video sslmode=disable",
	}
}

func (h *Hp) String() string {
	return h.Host + ":" + strconv.Itoa(h.Port)
}

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
