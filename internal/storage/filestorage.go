package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

type shortenUrlData struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewRead(filename string, strg *Storage) {
	//filePath := filepath.Join(directory, filename)
	//
	//if _, err := os.Stat(filePath); err == nil {
	//	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0666)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return &storageFile{
	//		file: file,
	//	}, nil
	//} else if os.IsNotExist(err) {
	//	files, err := os.ReadDir(directory)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	jsonFile := false
	//
	//	for _, fileInfo := range files {
	//		fmt.Println(fileInfo)
	//		if filepath.Ext(fileInfo.Name()) == ".json" {
	//			jsonFile = true
	//			jsFilePath := filepath.Join(directory, fileInfo.Name())
	//
	//			jsData, err := os.ReadFile(jsFilePath)
	//			if err != nil {
	//				return nil, err
	//			}
	//
	//			// Удаляем .js файл
	//			err = os.Remove(jsFilePath)
	//			if err != nil {
	//				return nil, err
	//			}
	//
	//			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//			if err != nil {
	//				return nil, err
	//			}
	//
	//			_, err = file.Write(jsData)
	//			if err != nil {
	//				return nil, err
	//			}
	//
	//			return &storageFile{
	//				file: file,
	//			}, nil
	//		}
	//		if !jsonFile {
	//			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//			if err != nil {
	//				return nil, err
	//			}
	//			return &storageFile{
	//				file: file,
	//			}, nil
	//		}
	//	}
	//} else {
	//	return nil, err
	//}
	//
	//return nil, nil

	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		var myData shortenUrlData
		if err := json.NewDecoder(strings.NewReader(line)).Decode(&myData); err == nil {
			strg.SetData(myData.OriginalURL, myData.OriginalURL)
		} else {
			log.Fatal(err)
		}
	}
}

func NewWrite(filename, link, shortlink string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	myData := shortenUrlData{
		ShortURL:    shortlink,
		OriginalURL: link,
	}

	if err := json.NewEncoder(file).Encode(myData); err != nil {
		log.Fatal(err)
	}
}
