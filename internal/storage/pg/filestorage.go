package pgstorage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"practicumserver/internal/models"
	"strings"
)

// Comments for the GetData, SetData, SetListData, GetListData, CheckUserID, UpdateDeletedFlag
// methods are in storage/StorageMock

// FileStorage structure for storing data in a file
type FileStorage struct {
	FileName string
	*MemoryStorage
	StorageParam
}

// NewRead method reads data from a file and transfers it to the fields of the MemoryStorage structure.
func (fs *FileStorage) NewRead() error {
	file, err := os.OpenFile(fs.FileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		var fileData models.FileData
		if err := json.NewDecoder(strings.NewReader(line)).Decode(&fileData); err == nil {
			fs.SetFromFileData(&fileData)
		} else {
			return err
		}
	}
	return nil
}

// FullWrite a method for overwriting data to a file when deleting it from a function UpdateDeletedFlag.
func (fs *FileStorage) FullWrite() error {
	file, err := os.OpenFile(fs.FileName, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	for UserID, URLList := range fs.UserIDUrls {
		for shortLink, originalURL := range URLList {
			URLData := models.FileData{
				UserID:      UserID,
				ShortURL:    shortLink,
				OriginalURL: originalURL,
				DeletedFlag: false,
			}
			if _, ok := fs.DeletedURL[shortLink]; ok {
				URLData.DeletedFlag = true
			}
			if err := json.NewEncoder(file).Encode(URLData); err != nil {
				return err
			}
		}
	}
	return nil
}

// NewWrite method writes data to a file.
func (fs *FileStorage) NewWrite(userIDStr, originalURL, ShortURL string) {
	file, err := os.OpenFile(fs.FileName, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileData := models.FileData{
		UserID:      userIDStr,
		ShortURL:    ShortURL,
		OriginalURL: originalURL,
		DeletedFlag: false,
	}

	if err := json.NewEncoder(file).Encode(fileData); err != nil {
		log.Fatal(err)
	}
}

// SetFromFileData method sets data from a file into the MemoryStorage fields.
func (fs *FileStorage) SetFromFileData(fileData *models.FileData) {
	fs.LinkBoolUrls[fileData.OriginalURL] = true
	fs.ShortBoolUrls[fileData.ShortURL] = false
	if !fileData.DeletedFlag {
		fs.ShortUrls[fileData.ShortURL] = fileData.OriginalURL
	} else {
		fs.DeletedURL[fileData.ShortURL] = fileData.OriginalURL
	}
	if fs.UserIDUrls[fileData.UserID] == nil {
		fs.UserIDUrls[fileData.UserID] = make(map[string]string)
	}
	fs.UserIDUrls[fileData.UserID][fileData.ShortURL] = fileData.OriginalURL
}

func (fs *FileStorage) SetData(ctx context.Context, originalURL string) (string, error) {
	if _, ok := fs.LinkBoolUrls[originalURL]; !ok {
		shortLink, err := fs.MemoryStorage.SetData(ctx, originalURL)
		if err != nil {
			return "", err
		}
		if userIDStr, ok := ctx.Value(models.UserIDKey).(string); ok {
			fs.NewWrite(userIDStr, originalURL, shortLink)
			return shortLink, nil
		}
		return "", errors.New("UserID is not valid type string")
	}
	return fs.checkShortLink(originalURL)
}

func (fs *FileStorage) SetListData(ctx context.Context,
	reqList []models.RequestAPIBatch, prefix string) ([]models.ResponseAPIBatch, error) {

	respList := make([]models.ResponseAPIBatch, 0)

	if userIDStr, ok := ctx.Value(models.UserIDKey).(string); ok {
		for _, structOriginalURL := range reqList {
			if _, ok := fs.LinkBoolUrls[structOriginalURL.OriginalURL]; !ok {
				shortLink, err := fs.MemoryStorage.SetData(ctx, structOriginalURL.OriginalURL)
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
					fs.NewWrite(userIDStr, structOriginalURL.OriginalURL, shortLink)
				}
			} else {
				shortLink, _ := fs.checkShortLink(structOriginalURL.OriginalURL)
				resp := models.ResponseAPIBatch{
					CorrelationID: structOriginalURL.CorrelationID,
					ShortURL:      prefix + "/" + shortLink,
				}
				respList = append(respList, resp)
			}
		}
		return respList, nil
	}
	return nil, errors.New("UserID is not valid type string")
}

func (fs *FileStorage) GetListData(ctx context.Context, prefix string) ([]models.ResponseAPIUserUrls, error) {
	return fs.MemoryStorage.GetListData(ctx, prefix)
}

func (fs *FileStorage) UpdateDeletedFlag(ctx context.Context, userIDList, shortLinkList []string) error {
	if err := fs.MemoryStorage.UpdateDeletedFlag(ctx, userIDList, shortLinkList); err != nil {
		return err
	}
	if err := fs.FullWrite(); err != nil {
		return err
	}
	return nil
}
