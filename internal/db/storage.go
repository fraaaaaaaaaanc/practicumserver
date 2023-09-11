package db

import (
	"errors"
)

var shortUrls = map[string]string{
	"http://test": "test", // Значение "http://test" заданно для тестирования
}

func SetDB(key, id string) (string, error) {
	if _, ok := shortUrls[key]; !ok {
		shortUrls[key] = id
		return id, nil
	}
	return shortUrls[key], errors.New("key already exists")
}

func GetDB(link string) (string, error) {
	for k, v := range shortUrls {
		if v == link {
			return k, nil
		}
	}
	return "", errors.New("the initial link is missing")
}
