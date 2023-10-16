package storage

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgconn"
	"practicumserver/internal/models"
	"practicumserver/internal/utils"
	"time"
)

func (ds *DBStorage) PingDB(ctx context.Context) error {
	if err := ds.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (ds *DBStorage) checkShortLink(ctx context.Context) (string, error) {
	for {
		shortLink := utils.LinkShortening()

		var exists bool
		row := ds.db.QueryRowContext(ctx,
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

func (ds *DBStorage) getNewShortLink(ctx context.Context, link string) (string, error) {
	var shortlink string
	row := ds.db.QueryRowContext(ctx,
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
	row := ds.db.QueryRowContext(ctx,
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
		_, err = ds.db.ExecContext(ctx,
			"INSERT INTO links (Link, ShortLink) "+
				"VALUES ($1, $2)",
			originalURL, shortLink)
		if err != nil {
			if pqErr, ok := err.(*pgconn.PgError); ok {
				if pqErr.Code == "23505" {
					return shortLink, nil
				}
			}
			return "", err
		}
		return shortLink, nil
	}
}

func (ds *DBStorage) SetListData(ctx context.Context,
	reqList []models.RequestAPIBatch, prefix string) ([]models.ResponseAPIBatch, error) {
	ds.sm.Lock()
	defer ds.sm.Unlock()

	tx, err := ds.db.Begin()
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
				"INSERT INTO links (Link, ShortLink) "+
					"VALUES ($1, $2)",
				StructOriginalURL.OriginalURL, shortLink)
			if err != nil {
				if pqErr, ok := err.(*pgconn.PgError); ok {
					if pqErr.Code != "23505" {
						tx.Rollback()
						return nil, err
					}
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
