package storage

type StorageMock interface {
	CheckShortLink() (string, error)
	GetNewShortLink(link string) (string, error)
	SetData(link, shortLink string) error
	GetData(shortLink string) (string, error)
}
