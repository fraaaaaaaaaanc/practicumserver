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

type shortenURLData struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

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

func (fs *FileStorage) SetFromFileData(originalURL, shortLink string) {
	fs.LinkBoolUrls[originalURL] = true
	fs.ShortBoolUrls[shortLink] = false
	fs.ShortUrls[shortLink] = originalURL
}

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
	return fs.checkShortLink(originalURL), nil
}

func (fs *FileStorage) SetListData(ctx context.Context,
	reqList []models.RequestAPIBatch) ([]models.ResponseAPIBatch, error) {

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
					ShortURL:      shortLink,
				}
				respList = append(respList, resp)
				fs.NewWrite(structOriginalURL.OriginalURL, shortLink)
			}
		} else {
			resp := models.ResponseAPIBatch{
				CorrelationID: structOriginalURL.CorrelationID,
				ShortURL:      fs.checkShortLink(structOriginalURL.OriginalURL),
			}
			respList = append(respList, resp)
		}
	}
	return respList, nil

	//respList, err := fs.MemoryStorage.SetListData(ctx, reqList)
	//if err != nil {
	//	return nil, err
	//}
	//select {
	//case <-ctx.Done():
	//	return nil, err
	//default:
	//	for idx, structOriginalUrl := range reqList {
	//		if _, ok := fs.LinkBoolUrls[structOriginalUrl.OriginalURL]; !ok {
	//			fs.NewWrite(structOriginalUrl.OriginalURL, respList[idx].ShortURL)
	//		}
	//	}
	//}
	//return respList, nil
}
