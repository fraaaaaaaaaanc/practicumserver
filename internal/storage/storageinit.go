package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"practicumserver/internal/storage/pg"
	"sync"
)

// NewStorage creates and returns a storage instance based on the provided parameters.
// It accepts a logger for logging, a DBStorageAdress for database storage, and a FileStoragePath for file-based storage.
// The function determines the storage type based on the presence of these parameters.
func NewStorage(log *zap.Logger,
	DBStorageAdress, FileStoragePath string) (StorageMock, error) {
	var sm sync.Mutex
	// Create a storage structure with common elements for each storage type.
	strg := pgstorage.StorageParam{
		Log: log,
		Sm:  &sm,
	}
	// Create a storage for working with a database.
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
		//Данный запрос создает таблицу links в БД, если ее там нет,
		//Данная таблица имеет три поля id, Link, ShortLink
		//Поле id является PRIMARY KEY, а поля Link, ShortLink не могу повторятся (UNIQUE)
		_, err = db.ExecContext(ctx, `
			DO $$ 
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'links') THEN
					CREATE TABLE links (
						id SERIAL PRIMARY KEY, 
						UserID VARCHAR(24) NOT NULL DEFAULT '', 
						Link VARCHAR(250) NOT NULL DEFAULT '' UNIQUE,
						ShortLink VARCHAR(250) NOT NULL DEFAULT '' UNIQUE,
						DeletedFlag BOOLEAN NOT NULL DEFAULT false
					);
				END IF;
			END $$;
			`)
		if err != nil {
			return nil, err
		}
		// Return the created database storage.
		return &pgstorage.DBStorage{
			DB:           db,
			StorageParam: strg,
		}, nil
	}
	// Create a memory storage.
	// This storage is created before the file-based storage since it's used as an anonymous field in the latter.
	// The memory storage initializes some map fields for testing.
	memoryStorage := &pgstorage.MemoryStorage{
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
		UserIDUrls: map[string]map[string]string{
			"test": {"test": "http://test"},
		},
		DeletedURL: map[string]string{
			"test": "http://test",
		},
	}
	// Create a file storage instance and read data from the file.
	if FileStoragePath != "" {
		fs := &pgstorage.FileStorage{
			MemoryStorage: memoryStorage,
			FileName:      FileStoragePath,
		}
		if err := fs.NewRead(); err != nil {
			return nil, err
		}
		return fs, nil
	}
	return memoryStorage, nil
}
