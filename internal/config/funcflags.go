package config

import (
	"errors"
	"strconv"
	"strings"
)

// Струкутра описывающая данные адреса сервера
type Hp struct {
	Host string
	Port int
}

// Структура описывающая поля для хранения данных флагов
type Flags struct {
	Prefix          string
	LogLevel        string
	FileLog         bool
	ConsoleLog      bool
	FileStoragePath string
	DBStorageAdress string
	Hp
}

// Инициализатор структры Flags задающий значения по умолчанию
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

// Переопределение метода String для формирования адресса сервера
func (h *Hp) String() string {
	return h.Host + ":" + strconv.Itoa(h.Port)
}

// Переопределение метода Set для проверки адресса сервера и записи данных в структуру Hp
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
