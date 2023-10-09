package storage

type StorageMock interface {
	СheckShortLink(filename, link string) string
	GetNewShortLink(link, filename string) string
	SetData(link, shortLink string)
	GetData(shortLink string) (string, bool)
}
