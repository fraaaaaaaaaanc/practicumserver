package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"practicumserver/internal/models"
	"strings"
)

// Структура для хранения данных при их чтении из файла методом NewRead
type shortenURLData struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

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

		var myData shortenURLData
		if err := json.NewDecoder(strings.NewReader(line)).Decode(&myData); err == nil {
			fs.SetFromFileData(myData.OriginalURL, myData.ShortURL)
		} else {
			return err
		}
	}
	return nil
}

// Метод для записи данных в файл
func (fs *FileStorage) NewWrite(originalURL, ShortURL string) {
	file, err := os.OpenFile(fs.FileName, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	myData := shortenURLData{
		ShortURL:    ShortURL,
		OriginalURL: originalURL,
	}

	if err := json.NewEncoder(file).Encode(myData); err != nil {
		log.Fatal(err)
	}
}

// Метод для записи данных в поля структура MemoryStorage при чтении их из файла
func (fs *FileStorage) SetFromFileData(originalURL, shortLink string) {
	fs.LinkBoolUrls[originalURL] = true
	fs.ShortBoolUrls[shortLink] = false
	fs.ShortUrls[shortLink] = originalURL
}

// Переопределение метожа SetData структуры MemoryStorage
// Метод вызывает SetData после чего записывает данные в файл методом NewWrite
func (fs *FileStorage) SetData(ctx context.Context, originalURL string) (string, error) {
	if _, ok := fs.LinkBoolUrls[originalURL]; !ok {
		shortLink, err := fs.MemoryStorage.SetData(ctx, originalURL)
		if err != nil {
			return "", err
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			fs.NewWrite(originalURL, shortLink)
			return shortLink, nil
		}
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
				fs.NewWrite(structOriginalURL.OriginalURL, shortLink)
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
