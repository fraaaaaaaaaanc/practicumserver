package storage

import (
	"practicumserver/internal/utils"
	"sync"
)

type Storage struct {
	ShortBoolUrls map[string]bool
	LinkBoolUrls  map[string]bool
	ShortUrls     map[string]string
	sm            sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		ShortBoolUrls: map[string]bool{
			"test": true,
		},
		LinkBoolUrls: map[string]bool{
			"http://test": true,
		},
		ShortUrls: map[string]string{
			"test": "http://test",
		},
	}
}

func (s *Storage) СheckShortLink() string {
	shortLink := utils.LinkShortening()
	for s.ShortBoolUrls[shortLink] {
		shortLink = utils.LinkShortening()
	}
	return shortLink
}

func (s *Storage) GetNewShortLink(link string) string {
	s.sm.Lock()
	defer s.sm.Unlock()
	if _, ok := s.LinkBoolUrls[link]; ok {
		for shortLink, longLink := range s.ShortUrls {
			if longLink == link {
				return shortLink
			}
		}
	}
	return s.СheckShortLink()
}

func (s *Storage) SetData(link, shortLink string) {
	s.sm.Lock()
	defer s.sm.Unlock()
	s.ShortUrls[shortLink] = link
	s.ShortBoolUrls[shortLink] = true
	s.LinkBoolUrls[link] = true
}

func (s *Storage) GetData(shortLink string) (string, bool) {
	s.sm.Lock()
	defer s.sm.Unlock()
	if _, ok := s.ShortBoolUrls[shortLink]; ok {
		return s.ShortUrls[shortLink], false
	}
	return "", true
}
