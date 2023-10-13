package storage

import (
	"context"
	"practicumserver/internal/utils"
)

func (ms *MemoryStorage) CheckShortLink(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		shortLink := utils.LinkShortening()
		for ms.ShortBoolUrls[shortLink] {
			shortLink = utils.LinkShortening()
		}
		return shortLink, nil
	}
}

func (ms *MemoryStorage) GetNewShortLink(ctx context.Context, link string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		ms.sm.Lock()
		defer ms.sm.Unlock()
		if _, ok := ms.LinkBoolUrls[link]; ok {
			for shortLink, longLink := range ms.ShortUrls {
				if longLink == link {
					return shortLink, nil
				}
			}
		}
		shortLink, err := ms.CheckShortLink(ctx)
		if err != nil {
			return "", err
		}
		return shortLink, nil
	}
}

func (ms *MemoryStorage) GetData(ctx context.Context, shortLink string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		ms.sm.Lock()
		defer ms.sm.Unlock()
		if _, ok := ms.ShortBoolUrls[shortLink]; ok {
			return ms.ShortUrls[shortLink], nil
		}
		return "", nil
	}
}

func (ms *MemoryStorage) SetData(ctx context.Context, originalURL, shortLink string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		ms.sm.Lock()
		defer ms.sm.Unlock()
		if _, ok := ms.LinkBoolUrls[originalURL]; !ok {
			ms.ShortUrls[shortLink] = originalURL
			ms.ShortBoolUrls[shortLink] = true
			ms.LinkBoolUrls[originalURL] = true
		}
		return nil
	}
}
