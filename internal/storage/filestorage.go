package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type shortenURLData struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewRead(filename string, strg StorageMock) error {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		var myData shortenURLData
		if err := json.NewDecoder(strings.NewReader(line)).Decode(&myData); err == nil {
			strg.SetData(myData.OriginalURL, myData.OriginalURL)
		} else {
			return err
		}
	}
	return nil
}

func NewWrite(filename, link, shortlink string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("((((((")
		log.Fatal(err)
	}
	myData := shortenURLData{
		ShortURL:    shortlink,
		OriginalURL: link,
	}

	if err := json.NewEncoder(file).Encode(myData); err != nil {
		log.Fatal(err)
	}
}
