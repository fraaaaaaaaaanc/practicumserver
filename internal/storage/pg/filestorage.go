package pgstorage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"practicumserver/internal/models"
	"strings"
)

// Структура для хранения данных при их чтении из файла методом NewRead

// Структура для хранения данных в файле
type FileStorage struct {
	FileName string
	*MemoryStorage
	StorageParam
}

// Метод для чтениях данных из файла и их переноса в поля структура MemoryStorage
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
			if _, ok := fs.DeletedURl[shortLink]; ok {
				URLData.DeletedFlag = true
			}
			if err := json.NewEncoder(file).Encode(URLData); err != nil {
				return err
			}
		}
	}
	return nil
}

// Метод для записи данных в файл
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

// Метод для записи данных в поля структура MemoryStorage при чтении их из файла
func (fs *FileStorage) SetFromFileData(fileData *models.FileData) {
	fmt.Println(fileData)
	fs.LinkBoolUrls[fileData.OriginalURL] = true
	fs.ShortBoolUrls[fileData.ShortURL] = false
	if !fileData.DeletedFlag {
		fs.ShortUrls[fileData.ShortURL] = fileData.OriginalURL
	} else {
		fs.DeletedURl[fileData.ShortURL] = fileData.OriginalURL
	}
	if fs.UserIDUrls[fileData.UserID] == nil {
		fs.UserIDUrls[fileData.UserID] = make(map[string]string)
	}
	fs.UserIDUrls[fileData.UserID][fileData.ShortURL] = fileData.OriginalURL
}

// Переопределение метожа SetData структуры MemoryStorage
// Метод вызывает SetData после чего записывает данные в файл методом NewWrite
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

// Метод принимает слайс оригинальных URL reqList []models.RequestAPIBatch
// и проверяет их на наличие в этом хранилище, если данные на записаны то программа вызовет
// метод SetData и запишет их в результирующий слайс respList []models.ResponseAPIBatch,
// иначе программа вызовет метод checkShortLink, получит сокращенные URL для переданного OriginalURL
// и запишет результат в respList []models.ResponseAPIBatch
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
