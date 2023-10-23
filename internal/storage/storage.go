package storage

import (
	"context"
	"errors"
	"practicumserver/internal/models"
)

type StorageMock interface {
	//Метод принимает префикс для сокращенной сылки, оригинальную ссылку и вызывает для нее метод checkShortLink,
	//далее функция метод проверяет есть ли данная ссылка в хранилище, если да, то не устанавливает
	//значения, если нет, то устанавливает, метод возвращает сокращенную ссылку для переданного originalURL и
	//объект типа error
	SetData(ctx context.Context, prefix, link string) (string, error)
	//Метода который принимает сокращенную сыылку и проверяет есть ли она в хранилище,
	//если такая сокращенная ссылка уже есть, то функция возвращает оригинальную ссылку,
	//иначе функция возвращает пустую строку
	GetData(ctx context.Context, shortLink string) (string, error)
	//Метод принимает слайс []models.RequestAPIBatch с множеством оригинальных url
	//и последовательно для кажого из них вызывает метод SetData, после чего помещает полученный
	//сокращенный URL в слайс respList []models.ResponseAPIBatch
	//Метод возвращает слайс []models.ResponseAPIBatch и объект типа error
	SetListData(ctx context.Context, reqList []models.RequestAPIBatch, prefix string) ([]models.ResponseAPIBatch, error)
	//Метод принимает контекст, в который момент инициализации пользователя помещается его ID, после чего в хранилищце
	//происходит поиск данных привязанных к этому ID
	//Метод возвращает готовый ответ типа []models.ResponseApiUserUrls
	GetListData(ctx context.Context) ([]models.ResponseApiUserUrls, error)
}

// Ошибка сообщающая о конфликте данных в хранилище
var ErrConflictData = errors.New("data conflict: the resulting url already exists in the storage")
