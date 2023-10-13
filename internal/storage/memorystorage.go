package storage

import (
	"practicumserver/internal/utils"
)

func (ms *MemoryStorage) CheckShortLink() (string, error) {
	shortLink := utils.LinkShortening()
	for ms.ShortBoolUrls[shortLink] {
		shortLink = utils.LinkShortening()
	}
	return shortLink, nil
}

func (ms *MemoryStorage) GetNewShortLink(link string) (string, error) {
	ms.sm.Lock()
	defer ms.sm.Unlock()
	if _, ok := ms.LinkBoolUrls[link]; ok {
		for shortLink, longLink := range ms.ShortUrls {
			if longLink == link {
				return shortLink, nil
			}
		}
	}
	shortLink, err := ms.CheckShortLink()
	if err != nil {
		return "", err
	}
	return shortLink, nil
}

func (ms *MemoryStorage) GetData(shortLink string) (string, error) {
	ms.sm.Lock()
	defer ms.sm.Unlock()
	if _, ok := ms.ShortBoolUrls[shortLink]; ok {
		return ms.ShortUrls[shortLink], nil
	}
	return "", nil
}

func (ms *MemoryStorage) SetData(originalURL, shortLink string) error {
	ms.sm.Lock()
	defer ms.sm.Unlock()
	if _, ok := ms.LinkBoolUrls[originalURL]; !ok {
		ms.ShortUrls[shortLink] = originalURL
		ms.ShortBoolUrls[shortLink] = true
		ms.LinkBoolUrls[originalURL] = true
	}
	return nil
}
