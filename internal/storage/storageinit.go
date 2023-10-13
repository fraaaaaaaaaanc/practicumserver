package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"sync"
)

type StorageParam struct {
	log *zap.Logger
	sm  sync.Mutex
}

type DBStorage struct {
	db *sql.DB
	StorageParam
}

type FileStorage struct {
	FileName string
	*MemoryStorage
	StorageParam
}

type MemoryStorage struct {
	ShortBoolUrls map[string]bool
	LinkBoolUrls  map[string]bool
	ShortUrls     map[string]string
	StorageParam
}

func NewStorage(log *zap.Logger,
	DBStorageAdress, FileStoragePath string) (StorageMock, error) {
	var sm sync.Mutex
	strg := StorageParam{
		log: log,
		sm:  sm,
	}
	if DBStorageAdress != "" {
		db, err := sql.Open("pgx",
			DBStorageAdress)
		if err != nil {
			log.Error("Error:", zap.Error(err))
			return nil, err
		}

		ctx, cansel := context.WithCancel(context.Background())
		defer cansel()

		if err = db.PingContext(ctx); err != nil {
			log.Error("Error:", zap.Error(err))
			return nil, err
		}
		//row, err := db.QueryContext(ctx, "SELECT EXISTS "+
		//	"(SELECT 1 FROM information_schema.tables WHERE table_name = 'links')")
		//if err != nil {
		//	return nil, err
		//}
		//if row.Next() {
		//	err = row.Scan(&exists)
		//	if err != nil {
		//		return nil, err
		//	}
		//}
		//defer row.Close()
		//
		//if !exists {
		//	db.Exec("CREATE TABLE links (" +
		//		"id SERIAL PRIMARY KEY, " +
		//		"Link VARCHAR(250) NOT NULL DEFAULT '' UNIQUE," +
		//		"ShortLink VARCHAR(250) NOT NULL DEFAULT '' UNIQUE)")
		//}
		_, err = db.ExecContext(ctx, `
			DO $$ 
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'links') THEN
					CREATE TABLE links (
						id SERIAL PRIMARY KEY, 
						Link VARCHAR(250) NOT NULL DEFAULT '' UNIQUE,
						ShortLink VARCHAR(250) NOT NULL DEFAULT '' UNIQUE
					);
				END IF;
			END $$;
			`)
		if err != nil {
			return nil, err
		}
		return &DBStorage{
			db:           db,
			StorageParam: strg,
		}, nil
	}
	memoryStorage := &MemoryStorage{
		StorageParam: strg,
		ShortBoolUrls: map[string]bool{
			"test": true,
		},
		LinkBoolUrls: map[string]bool{
			"http://test": true,
		},
		ShortUrls: map[string]string{
			"test": "http://test",
		},
	}
	if FileStoragePath != "" {
		fs := &FileStorage{
			MemoryStorage: memoryStorage,
			FileName:      FileStoragePath,
		}
		fs.NewRead()
		return fs, nil
	}
	return memoryStorage, nil
}
