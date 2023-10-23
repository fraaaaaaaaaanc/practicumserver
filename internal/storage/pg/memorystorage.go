package storage

import (
	"context"
	"practicumserver/internal/models"
	"practicumserver/internal/storage"
	"practicumserver/internal/utils"
)

//Комментарии для методов SetData, GetData, SetListData находятся в storage/StorageMock

// Структура для хранения данных в памяти
type MemoryStorage struct {
	ShortBoolUrls map[string]bool
	LinkBoolUrls  map[string]bool
	ShortUrls     map[string]string
	StorageParam
}

// Метод который формирует новую сокращенную ссылку и проверяет
// cуществует ли такая сокращенная ссылка, если она есть, то функция
// генерирует новую сокращенную ссылку пока она не будет уникальной
func (ms *MemoryStorage) getNewShortLink() string {
	shortLink := utils.LinkShortening()
	for ms.ShortBoolUrls[shortLink] {
		shortLink = utils.LinkShortening()
	}
	return shortLink
}

// Метод который проверят наличие оригинальной ссылки в хранилище,
// если переданная оригинальная ссылка уже есть, то код возвращает ее сокращенный
// варинт и ошибку storage.ErrConflictData, иначе вызывает метод getNewShortLink
func (ms *MemoryStorage) checkShortLink(prefix, originalURL string) (string, error) {
	if _, ok := ms.LinkBoolUrls[originalURL]; ok {
		for shortLink, longLink := range ms.ShortUrls {
			if longLink == originalURL {
				return shortLink, storage.ErrConflictData
			}
		}
	}
	shortLink := prefix + "/" + ms.getNewShortLink()
	return shortLink, nil
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

func (ms *MemoryStorage) SetData(ctx context.Context, prefix, originalURL string) (string, error) {
	ms.sm.Lock()
	defer ms.sm.Unlock()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		shortLink, err := ms.checkShortLink(prefix, originalURL)
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
		shortLink, err := ms.SetData(ctx, prefix, structOriginalURL.OriginalURL)
		if err != nil {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			resp := models.ResponseAPIBatch{
				CorrelationID: structOriginalURL.CorrelationID,
				ShortURL:      shortLink,
			}
			respList = append(respList, resp)
		}
	}
	return respList, nil
}

func (ms *MemoryStorage) GetListData(ctx context.Context) ([]models.ResponseApiUserUrls, error) {
	return nil, nil
}
