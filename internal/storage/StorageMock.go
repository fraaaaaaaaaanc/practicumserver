package storage

import (
	"context"
	"errors"
	"practicumserver/internal/models"
)

type StorageMock interface {
	SetData(ctx context.Context, link string) (string, error)
	GetData(ctx context.Context, shortLink string) (string, error)
	SetListData(ctx context.Context, reqList []models.RequestAPIBatch, prefix string) ([]models.ResponseAPIBatch, error)
}

var ErrConflictData = errors.New("data conflict: the resulting url already exists in the storage")
