package storage

import "context"

type StorageMock interface {
	CheckShortLink(ctx context.Context) (string, error)
	GetNewShortLink(ctx context.Context, link string) (string, error)
	SetData(ctx context.Context, link, shortLink string) error
	GetData(ctx context.Context, shortLink string) (string, error)
}
