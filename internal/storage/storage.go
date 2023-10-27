package storage

import (
	"context"
	"practicumserver/internal/models"
)

type StorageMock interface {
	//Метод принимает оригинальную ссылку и вызывает для нее метод checkShortLink,
	//далее функция метод проверяет есть ли данная ссылка в хранилище, если да, то не устанавливает
	//значения, если нет, то устанавливает, метод возвращает сокращенную ссылку для переданного originalURL и
	//объект типа error
	SetData(ctx context.Context, link string) (string, error)
	//Метода который принимает сокращенную сыылку и проверяет есть ли она в хранилище,
	//если такая сокращенная ссылка уже есть, то функция возвращает оригинальную ссылку,
	//иначе функция возвращает пустую строку
	GetData(ctx context.Context, shortLink string) (string, error)
	//Метод принимает слайс []models.RequestAPIBatch с множеством оригинальных url
	//и последовательно для кажого из них вызывает метод SetData, после чего помещает полученный
	//сокращенный URL в слайс respList []models.ResponseAPIBatch
	//Метод возвращает слайс []models.ResponseAPIBatch и объект типа error
	SetListData(ctx context.Context, reqList []models.RequestAPIBatch, prefix string) ([]models.ResponseAPIBatch, error)
	GetListData(ctx context.Context, prefix string) ([]models.ResponseAPIUserUrls, error)
	CheckUserID(ctx context.Context, userID string) (bool, error)
	UpdateDeletedFlag(ctx context.Context, userIDList, shortLinkList []string) error
}
