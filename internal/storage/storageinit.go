package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"practicumserver/internal/storage/pg"
	"sync"
)

// Функция NewStorage принимает параметры log(для логирования), и две булевые переменные
// полученные при парсинге флагов. Функция проверяет эти флаги, если тот или иной флга
// принимает значени true, то функция создает объект storage того типа, реализующий
// интерфейс storage.StorageMock
func NewStorage(log *zap.Logger,
	DBStorageAdress, FileStoragePath string) (StorageMock, error) {
	//Cоздание структуры с общими элементами для кажлого storage
	var sm sync.Mutex
	strg := pgstorage.StorageParam{
		Log: log,
		Sm:  &sm,
	}
	//Создание storage для работы с БД
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
		return &pgstorage.DBStorage{
			DB:           db,
			StorageParam: strg,
		}, nil
	}
	//Создание storage для хранения данных в памяти
	//Данный storage создается раньше чем storage для работы с файлами
	//т.к. он является анонимным полем второго
	//При создании этого storage сразу создаются некоторе поля в map-ах
	//для тестирования
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
	//Создание storage для хранения данных в файле
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
