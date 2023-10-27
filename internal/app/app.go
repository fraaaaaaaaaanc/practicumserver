package app

import (
	"context"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"practicumserver/internal/config"
	handlers "practicumserver/internal/handlers/allhandlers"
	"practicumserver/internal/logger"
	"practicumserver/internal/models"
	"practicumserver/internal/router"
	"practicumserver/internal/storage"
	"time"
)

type app struct {
	Flags  *config.Flags
	Log    *logger.ZapLogger
	Strg   storage.StorageMock
	Hndlrs *handlers.Handlers
	Rtr    chi.Router
}

func NewApp() (*app, error) {
	//Парсинг флагов и создание переменной логирования
	flags := config.ParseConfFlugs()
	log, err := logger.NewZapLogger(flags.FileLog)
	if err != nil {
		return nil, err
	}
	//Создание объекта storage реализующего интерфейсный тип storage.StorageMock
	strg, err := storage.NewStorage(log.Logger, flags.DBStorageAdress, flags.FileStoragePath)
	if err != nil {
		return nil, err
	}
	//Создание объекта handlers
	hndlrs := handlers.NewHandlers(strg, log.Logger, flags.Prefix)
	//Создание объекта роутера для передачи в http.ListenAndServe
	rtr, err := router.Router(hndlrs, log.Logger)
	if err != nil {
		return nil, err
	}

	appObj := &app{
		Flags:  flags,
		Log:    log,
		Strg:   strg,
		Hndlrs: hndlrs,
		Rtr:    rtr,
	}
	go appObj.DeleteShortLink()

	return appObj, nil
}

func (a *app) DeleteShortLink() {
	ticker := time.NewTicker(1 * time.Second)

	var shortLinkList []*models.DeleteURL
	for {
		select {
		case shortLinksStruct := <-a.Hndlrs.DelCn:
			shortLinkList = append(shortLinkList, shortLinksStruct)
		case <-ticker.C:
			if len(shortLinkList) == 0 {
				continue
			}
			go func() {
				var userID []string
				var shortLink []string

				for _, shortLinkStruct := range shortLinkList {
					userID = append(userID, shortLinkStruct.UserID)
					shortLink = append(shortLink, shortLinkStruct.ShortURL)
				}
				if err := a.Strg.UpdateDeletedFlag(context.Background(), userID, shortLink); err != nil {
					a.Log.Error("Error:", zap.Error(err))
				}
				shortLinkList = nil
			}()
		}
	}
}
