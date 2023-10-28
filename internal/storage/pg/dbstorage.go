package pgstorage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"practicumserver/internal/models"
	"practicumserver/internal/utils"
	"time"
)

// Comments for the GetData, SetData, SetListData, GetListData, CheckUserID, UpdateDeletedFlag
// methods are in storage/StorageMock

// DBStorage structure for storing data in a database
type DBStorage struct {
	DB *sql.DB
	StorageParam
}

// PingDB method checks the connection to the database.
func (ds *DBStorage) PingDB(ctx context.Context) error {
	if err := ds.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

// checkShortLink method generates a new short link and checks if it already exists in the database.
func (ds *DBStorage) checkShortLink(ctx context.Context) (string, error) {
	for {
		shortLink := utils.LinkShortening()

		var exists bool
		// This query checks if a record with shortLink exists in the ShortLink column and stores the
		// result in a boolean variable.
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

// getNewShortLink method checks for the existence of an original link in the storage.
// If the given original link already exists, it returns its shortened version and a storage.ErrConflictData error.
// Otherwise, it calls the getNewShortLink method.
func (ds *DBStorage) getNewShortLink(ctx context.Context, link string) (string, error) {
	var shortlink string
	// This query looks for a record with ShortLink for which Link = link (Original URL).
	// If the record is not found, the checkShortLink method is called.
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
	ctx, cansel := context.WithTimeout(ctx, 5*time.Second)
	defer cansel()

	var getResp GetResponse
	row := ds.DB.QueryRowContext(ctx,
		"SELECT Link, DeletedFlag FROM links WHERE ShortLink= $1",
		shortLink)

	if err := row.Scan(&getResp.originalURL, &getResp.deletedFlag); err != nil {
		if err == sql.ErrNoRows {
			return "", models.ErrNoRows
		}
		return "", err
	}
	if getResp.deletedFlag {
		return "", models.ErrDeletedData
	}
	return getResp.originalURL, nil
}

func (ds *DBStorage) SetData(ctx context.Context, originalURL string) (string, error) {
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
			ctx.Value(models.UserIDKey), originalURL, shortLink)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
				err = models.ErrConflictData
			}
			return shortLink, err
		}
		return shortLink, nil
	}
}

func (ds *DBStorage) SetListData(ctx context.Context,
	reqList []models.RequestAPIBatch, prefix string) ([]models.ResponseAPIBatch, error) {

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
				ctx.Value(models.UserIDKey), StructOriginalURL.OriginalURL, shortLink)
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
		ctx.Value(models.UserIDKey))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var oneResp models.ResponseAPIUserUrls
		if err = rows.Scan(&oneResp.ShortURL, &oneResp.OriginalURL); err != nil {
			return nil, err
		}
		oneResp.ShortURL = prefix + "/" + oneResp.ShortURL
		resp = append(resp, oneResp)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
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

func (ds *DBStorage) UpdateDeletedFlag(ctx context.Context, userIDList, shortLinkList []string) error {
	_, err := ds.DB.ExecContext(ctx,
		"UPDATE links Set DeletedFlag = true WHERE UserId = ANY ($1) AND ShortLink = ANY ($2)",
		userIDList, shortLinkList)
	if err != nil {
		return err
	}
	return nil
}
