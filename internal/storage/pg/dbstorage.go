package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"practicumserver/internal/models"
	"practicumserver/internal/storage"
	"practicumserver/internal/utils"
	"time"
)

//Комментарии для методов SetData, GetData, SetListData находятся в storage/StorageMock

// Структура для хранения данных в БД
type DBStorage struct {
	DB *sql.DB
	StorageParam
}

// Метод проверяющий подключение к БД
func (ds *DBStorage) PingDB(ctx context.Context) error {
	if err := ds.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

// Метод который формирует новую сокращенную ссылку и проверяет
// cуществует ли такая сокращенная ссылка, если она есть, то функция
// генерирует новую сокращенную ссылку пока она не будет уникальной
func (ds *DBStorage) checkShortLink(ctx context.Context) (string, error) {
	for {
		shortLink := utils.LinkShortening()

		var exists bool
		//Данный запрос проверяет существует ли запись shortLink в колонке ShortLink
		//и помещает результат в булевую переменную
		row := ds.DB.QueryRowContext(ctx,
			"SELECT EXISTS (SELECT 1 FROM links WHERE ShortLink = $1)",
			shortLink)
		if err := row.Scan(&exists); err != nil {
			return "", err
		}
		if !exists {
			return shortLink, nil
		}
	}
}

// Метод который проверят наличие оригинальной ссылки в хранилище,
// если переданная оригинальная ссылка уже есть, то код возвращает ее сокращенный
// варинт и ошибку storage.ErrConflictData, иначе вызывает метод getNewShortLink
func (ds *DBStorage) getNewShortLink(ctx context.Context, link string) (string, error) {
	var shortlink string
	//Данный запрос ищет запись ShortLink для которой Link = link(Оригинальная ссылка),
	//если запись не найдена вызывается метод checkShortLink
	row := ds.DB.QueryRowContext(ctx,
		"SELECT ShortLink FROM links WHERE Link = $1",
		link)
	if err := row.Scan(&shortlink); err != nil {
		if err == sql.ErrNoRows {
			shortLink, err := ds.checkShortLink(ctx)
			if err != nil {
				return "", err
			}
			return shortLink, nil
		}
		return "", err
	}
	return shortlink, nil
}

func (ds *DBStorage) GetData(ctx context.Context, shortLink string) (string, error) {
	ds.sm.Lock()
	defer ds.sm.Unlock()

	ctx, cansel := context.WithTimeout(ctx, 5*time.Second)
	defer cansel()

	var originLink string
	row := ds.DB.QueryRowContext(ctx,
		"SELECT Link FROM links WHERE ShortLink= $1",
		shortLink)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		if err := row.Scan(&originLink); err != nil {
			if err == sql.ErrNoRows {
				return "", nil
			}
			return "", err
		}
		return originLink, nil
	}
}

func (ds *DBStorage) SetData(ctx context.Context, originalURL string) (string, error) {
	ds.sm.Lock()
	defer ds.sm.Unlock()
	ctx, cansel := context.WithTimeout(ctx, 5*time.Second)
	defer cansel()

	shortLink, err := ds.getNewShortLink(ctx, originalURL)
	if err != nil {
		return "", err
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		_, err = ds.DB.ExecContext(ctx,
			"INSERT INTO links (UserID, Link, ShortLink) "+
				"VALUES ($1, $2, $3)",
			ctx.Value("userID"), originalURL, shortLink)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
				err = storage.ErrConflictData
			}
			return shortLink, err
		}
		return shortLink, nil
	}
}

func (ds *DBStorage) SetListData(ctx context.Context,
	reqList []models.RequestAPIBatch, prefix string) ([]models.ResponseAPIBatch, error) {
	ds.sm.Lock()
	defer ds.sm.Unlock()

	tx, err := ds.DB.Begin()
	if err != nil {
		return nil, err
	}

	ctx, cansel := context.WithTimeout(ctx, 5*time.Second)
	defer cansel()
	respList := make([]models.ResponseAPIBatch, 0)

	for _, StructOriginalURL := range reqList {
		shortLink, err := ds.getNewShortLink(ctx, StructOriginalURL.OriginalURL)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		select {
		case <-ctx.Done():
			tx.Rollback()
			return nil, err
		default:
			_, err = tx.ExecContext(ctx,
				"INSERT INTO links (UserID, Link, ShortLink) "+
					"VALUES ($1, $2, $3)",
				ctx.Value("userID"), StructOriginalURL.OriginalURL, shortLink)
			if err != nil {
				var pqErr *pgconn.PgError
				if errors.As(err, &pqErr) && pgerrcode.UniqueViolation == pqErr.Code {
					tx.Rollback()
					return nil, err
				}
			}
			resp := models.ResponseAPIBatch{
				CorrelationID: StructOriginalURL.CorrelationID,
				ShortURL:      prefix + "/" + shortLink,
			}
			respList = append(respList, resp)
		}
	}
	tx.Commit()
	return respList, nil
}

func (ds *DBStorage) GetListData(ctx context.Context, prefix string) ([]models.ResponseAPIUserUrls, error) {
	var resp []models.ResponseAPIUserUrls
	rows, err := ds.DB.QueryContext(ctx,
		"SELECT ShortLink, Link FROM links WHERE UserID = $1",
		ctx.Value("userID"))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var oneResp models.ResponseAPIUserUrls
		if err = rows.Scan(&oneResp.ShortURL, &oneResp.OriginalURL); err != nil {
			return nil, err
		}
		oneResp.ShortURL = prefix + "/" + oneResp.ShortURL
		resp = append(resp, oneResp)
	}
	return resp, nil
}

func (ds *DBStorage) CheckUserID(ctx context.Context, userID string) (bool, error) {
	row := ds.DB.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM links WHERE UserID = $1)",
		userID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	if !exists {
		return true, nil
	}
	return false, nil
}
