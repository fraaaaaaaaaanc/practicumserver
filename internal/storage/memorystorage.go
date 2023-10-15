package storage

import (
	"context"
	"practicumserver/internal/models"
	"practicumserver/internal/utils"
)

func (ms *MemoryStorage) checkShortLink() string {
	shortLink := utils.LinkShortening()
	for ms.ShortBoolUrls[shortLink] {
		shortLink = utils.LinkShortening()
	}
	return shortLink
}

func (ms *MemoryStorage) getNewShortLink(link string) string {
	if _, ok := ms.LinkBoolUrls[link]; ok {
		for shortLink, longLink := range ms.ShortUrls {
			if longLink == link {
				return shortLink
			}
		}
	}
	return ms.checkShortLink()
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
		shortLink := ms.getNewShortLink(originalURL)
		if _, ok := ms.LinkBoolUrls[originalURL]; !ok {
			ms.ShortUrls[shortLink] = originalURL
			ms.ShortBoolUrls[shortLink] = true
			ms.LinkBoolUrls[originalURL] = true
		}
		return shortLink, nil
	}
}

func (ms *MemoryStorage) SetListData(ctx context.Context,
	reqList []models.RequestApiBatch) ([]models.ResponseApiBatch, error) {
	//ms.sm.Lock()
	//defer ms.sm.Unlock()

	respList := make([]models.ResponseApiBatch, 0)

	for _, structOriginalURL := range reqList {
		shortLink, err := ms.SetData(ctx, structOriginalURL.OriginalUrl)
		if err != nil {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			resp := models.ResponseApiBatch{
				CorrelationID: structOriginalURL.CorrelationID,
				ShortURL:      shortLink,
			}
			respList = append(respList, resp)
		}
	}
	return respList, nil
}
