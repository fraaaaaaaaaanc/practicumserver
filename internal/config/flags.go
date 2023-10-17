package config

import (
	"flag"
	"os"
)

// Метод запускающий парсинг флагов,
// сначала метод проверят флаги переданные при запускеметодом ParseFlags, после чего
// вызывает метод ParseEnv для проверки переменных окружения, после чего возвращает структуру со значениями
// флагов
func ParseConfFlugs() *Flags {
	flags := ParseFlags()
	ParseEnv(flags)
	return flags
}

// Метод проверяет наличие значений в переменных окружения, если переменные не пустые, то
// их значения помещается в поле структуры flags, тем самым заменяя значение переданное пользователем
// при запуске программы или значение по умолчанию
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

// Метод задающий флаги командой строки и запускающий их парсинг
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
