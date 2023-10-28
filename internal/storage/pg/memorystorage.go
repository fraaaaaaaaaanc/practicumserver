package pgstorage

import (
	"context"
	"errors"
	"practicumserver/internal/models"
	"practicumserver/internal/utils"
)

// Comments for the GetData, SetData, SetListData, GetListData, CheckUserID, UpdateDeletedFlag
// methods are in storage/StorageMock

// MemoryStorage structures for in-memory storage
type MemoryStorage struct {
	ShortBoolUrls map[string]bool
	LinkBoolUrls  map[string]bool
	ShortUrls     map[string]string
	UserIDUrls    map[string]map[string]string
	DeletedURL    map[string]string
	StorageParam
}

// getNewShortLink method for generating a new short link and checking for its uniqueness.
func (ms *MemoryStorage) getNewShortLink() string {
	shortLink := utils.LinkShortening()
	for ms.ShortBoolUrls[shortLink] {
		shortLink = utils.LinkShortening()
	}
	return shortLink
}

// checkShortLink checks for the existence of an original URL in the storage. If the given original URL already exists,
// the function returns its shortened version and an error models.ErrConflictData. Otherwise, it calls the getNewShortLink method.
func (ms *MemoryStorage) checkShortLink(originalURL string) (string, error) {
	if _, ok := ms.LinkBoolUrls[originalURL]; ok {
		for shortLink, longLink := range ms.ShortUrls {
			if longLink == originalURL {
				return shortLink, models.ErrConflictData
			}
		}
	}
	return ms.getNewShortLink(), nil
}

func (ms *MemoryStorage) GetData(ctx context.Context, shortLink string) (string, error) {
	ms.Sm.Lock()
	defer ms.Sm.Unlock()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		if _, ok := ms.ShortUrls[shortLink]; ok {
			return ms.ShortUrls[shortLink], nil
		} else if _, ok := ms.DeletedURL[shortLink]; ok {
			return "", models.ErrDeletedData
		}
		return "", models.ErrNoRows
	}
}

func (ms *MemoryStorage) SetData(ctx context.Context, originalURL string) (string, error) {
	ms.Sm.Lock()
	defer ms.Sm.Unlock()

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

func (ms *MemoryStorage) UpdateDeletedFlag(ctx context.Context, userIDList, shortLinkList []string) error {
	var idx int
	for _, shortLink := range shortLinkList {
		if _, ok := ms.UserIDUrls[userIDList[idx]][shortLink]; ok {
			ms.DeletedURL[shortLink] = ms.ShortUrls[shortLink]
			delete(ms.ShortUrls, shortLink)
		}
		idx++
	}
	return nil
}
