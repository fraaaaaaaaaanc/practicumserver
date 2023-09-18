package storage

import (
	"errors"
	"sync"
)

type Storage struct {
	ShortBoolUrls map[string]bool
	ShortUrls     map[string]string
}

func (s *Storage) SetData(link, shortLink string) (string, error) {
	var sm sync.Mutex
	sm.Lock()
	defer sm.Unlock()
	if _, ok := s.ShortBoolUrls[shortLink]; !ok {
		s.ShortUrls[shortLink] = link
		s.ShortBoolUrls[shortLink] = true
		return link, nil
	}
	return s.ShortUrls[shortLink], errors.New("key already exists")
}

func (s *Storage) GetData(shortLink string) (string, error) {
	var sm sync.Mutex
	sm.Lock()
	defer sm.Unlock()
	if _, ok := s.ShortBoolUrls[shortLink]; ok {
		return s.ShortUrls[shortLink], nil
	}
	return "", errors.New("the initial link is missing")
}
