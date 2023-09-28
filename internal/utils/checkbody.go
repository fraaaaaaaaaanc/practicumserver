package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

func IsRequestBodyEmpty(body io.Reader) bool {
	fmt.Println(body)
	bodyCopy, err := ioutil.ReadAll(body)
	if err != nil {
		return true
	}

	// Создаем новый io.Reader на основе копии тела запроса
	bodyReader := bytes.NewReader(bodyCopy)

	// После создания копии и нового io.Reader, вы можете использовать bodyReader для проверки на пустоту
	if _, err := bodyReader.ReadByte(); err == io.EOF {
		return true
	}

	// Теперь bodyCopy содержит копию исходного тела запроса, и bodyReader можно использовать далее
	fmt.Println(bodyReader)
	return false
}
