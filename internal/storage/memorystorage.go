package storage

import (
	"context"
	"practicumserver/internal/models"
	"practicumserver/internal/utils"
)

func (ms *MemoryStorage) getNewShortLink() string {
	shortLink := utils.LinkShortening()
	for ms.ShortBoolUrls[shortLink] {
		shortLink = utils.LinkShortening()
	}
	return shortLink
}

func (ms *MemoryStorage) checkShortLink(originalURL string) (string, error) {
	if _, ok := ms.LinkBoolUrls[originalURL]; ok {
		for shortLink, longLink := range ms.ShortUrls {
			if longLink == originalURL {
				return shortLink, ErrConflictData
			}
		}
	}
	return ms.getNewShortLink(), nil
}

func (ms *MemoryStorage) GetData(ctx context.Context, shortLink string) (string, error) {
	ms.sm.Lock()
	defer ms.sm.Unlock()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		if _, ok := ms.ShortBoolUrls[shortLink]; ok {
			return ms.ShortUrls[shortLink], nil
		}
		return "", nil
	}
}

func (ms *MemoryStorage) SetData(ctx context.Context, originalURL string) (string, error) {
	ms.sm.Lock()
	defer ms.sm.Unlock()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		shortLink, err := ms.checkShortLink(originalURL)
		if _, ok := ms.LinkBoolUrls[originalURL]; !ok {
			ms.ShortUrls[shortLink] = originalURL
			ms.ShortBoolUrls[shortLink] = true
			ms.LinkBoolUrls[originalURL] = true
		}
		return shortLink, err
	}
}

func (ms *MemoryStorage) SetListData(ctx context.Context,
	reqList []models.RequestAPIBatch, prefix string) ([]models.ResponseAPIBatch, error) {
	//ms.sm.Lock()
	//defer ms.sm.Unlock()

	respList := make([]models.ResponseAPIBatch, 0)

	for _, structOriginalURL := range reqList {
		shortLink, err := ms.SetData(ctx, structOriginalURL.OriginalURL)
		if err != nil {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			resp := models.ResponseAPIBatch{
				CorrelationID: structOriginalURL.CorrelationID,
				ShortURL:      prefix + "/" + shortLink,
			}
			respList = append(respList, resp)
		}
	}
	return respList, nil
}
