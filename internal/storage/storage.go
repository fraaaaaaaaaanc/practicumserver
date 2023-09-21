package storage

import (
	"errors"
	"sync"
)

type Storage struct {
	ShortBoolUrls map[string]bool
	ShortUrls     map[string]string
	sm            sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		ShortBoolUrls: map[string]bool{
			"test": true,
		},
		ShortUrls: map[string]string{
			"test": "http://test",
		},
	}
}

func (s *Storage) SetData(link, shortLink string) (string, error) {
	s.sm.Lock()
	defer s.sm.Unlock()
	if _, ok := s.ShortBoolUrls[shortLink]; !ok {
		s.ShortUrls[shortLink] = link
		s.ShortBoolUrls[shortLink] = true
		return link, nil
	}
	return s.ShortUrls[shortLink], errors.New("key already exists")
}

func (s *Storage) GetData(shortLink string) (string, error) {
	s.sm.Lock()
	defer s.sm.Unlock()
	if _, ok := s.ShortBoolUrls[shortLink]; ok {
		return s.ShortUrls[shortLink], nil
	}
	return "", errors.New("the initial link is missing")
}
