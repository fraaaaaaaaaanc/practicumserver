package storage

import (
	"context"
	"errors"
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
	UserIDUrls    map[string]map[string]string
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
func (ms *MemoryStorage) checkShortLink(originalURL string) (string, error) {
	if _, ok := ms.LinkBoolUrls[originalURL]; ok {
		for shortLink, longLink := range ms.ShortUrls {
			if longLink == originalURL {
				return shortLink, storage.ErrConflictData
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

	shortLink, err := ms.checkShortLink(originalURL)
	if err != nil {
		return shortLink, err
	}
	userID := ctx.Value(models.UserIDKey)
	if userIDStr, ok := userID.(string); ok {
		if _, ok = ms.LinkBoolUrls[originalURL]; !ok {
			ms.ShortUrls[shortLink] = originalURL
			if ms.UserIDUrls[userIDStr] == nil {
				ms.UserIDUrls[userIDStr] = make(map[string]string)
			}
			ms.UserIDUrls[userIDStr][shortLink] = originalURL
			ms.ShortBoolUrls[shortLink] = true
			ms.LinkBoolUrls[originalURL] = true
			return shortLink, nil
		}
	}
	return "", errors.New("UserID is not valid type string")
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
		resp := models.ResponseAPIBatch{
			CorrelationID: structOriginalURL.CorrelationID,
			ShortURL:      prefix + "/" + shortLink,
		}
		respList = append(respList, resp)
	}
	return respList, nil
}

func (ms *MemoryStorage) GetListData(ctx context.Context, prefix string) ([]models.ResponseAPIUserUrls, error) {
	var resp []models.ResponseAPIUserUrls
	userID := ctx.Value(models.UserIDKey)
	if userIDStr, ok := userID.(string); ok {
		for key, elem := range ms.UserIDUrls[userIDStr] {
			oneResp := models.ResponseAPIUserUrls{
				ShortURL:    prefix + "/" + key,
				OriginalURL: elem,
			}
			resp = append(resp, oneResp)
		}
	}
	return resp, nil
}

func (ms *MemoryStorage) CheckUserID(ctx context.Context, userID string) (bool, error) {
	if _, ok := ms.UserIDUrls[userID]; !ok {
		return true, nil
	}
	return false, nil
}
