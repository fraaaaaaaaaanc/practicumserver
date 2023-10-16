package storage

import (
	"context"
	"practicumserver/internal/models"
)

type StorageMock interface {
	SetData(ctx context.Context, link string) (string, error)
	GetData(ctx context.Context, shortLink string) (string, error)
	SetListData(ctx context.Context, reqList []models.RequestAPIBatch, prefix string) ([]models.ResponseAPIBatch, error)
}
