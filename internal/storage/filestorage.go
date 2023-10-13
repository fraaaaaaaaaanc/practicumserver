package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
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
			fs.MemoryStorage.SetData(context.Background(), myData.OriginalURL, myData.ShortURL)
		} else {
			return err
		}
	}
	return nil
}

func (fs *FileStorage) NewWrite(originalURL, ShortURL string) {
	fs.sm.Lock()
	defer fs.sm.Unlock()
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

func (fs *FileStorage) SetData(ctx context.Context, originalURL, shortLink string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if _, ok := fs.LinkBoolUrls[originalURL]; !ok {
			fs.MemoryStorage.SetData(ctx, originalURL, shortLink)
			fs.NewWrite(originalURL, shortLink)
		}
		return nil
	}
}
