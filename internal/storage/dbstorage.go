package storage

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgconn"
	"practicumserver/internal/utils"
)

func (ds *DBStorage) PingDB(ctx context.Context) error {
	if err := ds.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (ds *DBStorage) CheckShortLink(ctx context.Context) (string, error) {
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

func (ds *DBStorage) GetNewShortLink(ctx context.Context, link string) (string, error) {
	ds.sm.Lock()
	defer ds.sm.Unlock()

	var shortlink string
	row := ds.db.QueryRowContext(ctx,
		"SELECT ShortLink FROM links WHERE Link = $1",
		link)
	if err := row.Scan(&shortlink); err != nil {
		if err == sql.ErrNoRows {
			shortLink, err := ds.CheckShortLink(ctx)
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

	var originLink string
	row := ds.db.QueryRowContext(ctx,
		"SELECT Link FROM links WHERE ShortLink= $1",
		shortLink)
	if err := row.Scan(&originLink); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return originLink, nil
}

func (ds *DBStorage) SetData(ctx context.Context, link, shortLink string) error {
	ds.sm.Lock()
	defer ds.sm.Unlock()
	//var boolOriginLink bool
	//row, err := ds.db.QueryContext(ds.ctx,
	//	"SELECT EXISTS (SELECT Link FROM links WHERE Link= $1)",
	//	link)
	//if err != nil {
	//	return err
	//}
	//if row.Next() {
	//	if err = row.Scan(&boolOriginLink); err != nil {
	//		return err
	//	}
	//}
	//fmt.Println("boolOriginLink", boolOriginLink)
	//if !boolOriginLink {
	//	_, err := ds.db.ExecContext(ds.ctx,
	//		"INSERT INTO links (Link, ShortLink) "+
	//			"VALUES ($1, $2)",
	//		link, shortLink)
	//	if err != nil {
	//		return err
	//	}
	//}
	//return nil
	_, err := ds.db.ExecContext(ctx,
		"INSERT INTO links (Link, ShortLink) "+
			"VALUES ($1, $2)",
		link, shortLink)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				return nil
			}
		}
		return err
	}
	return nil
}