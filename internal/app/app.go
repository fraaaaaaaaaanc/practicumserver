// Package app provides the main application logic, including initialization, configuration, and server setup.
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

// app represents the main application structure containing various components.// app represents the main application structure containing various components.
type app struct {
	Flags  *config.Flags       // Flags is a reference to the configuration flags used by the application.
	Log    *logger.ZapLogger   // Log is the application's logger, providing structured and efficient logging.
	Strg   storage.StorageMock // Strg represents the storage implementation used by the application.
	Hndlrs *handlers.Handlers  // Hndlrs is a reference to the application's handlers, providing request handling logic.
	Rtr    chi.Router          // Rtr is the application's router, responsible for routing incoming HTTP requests.
}

// New App initializes the application components and returns an app instance.
func NewApp() (*app, error) {
	// Parse command-line flags and create a logging variable.
	flags := config.ParseConfFlags()
	log, err := logger.NewZapLogger(flags.FileLog)
	if err != nil {
		return nil, err
	}
	//Создание объекта storage реализующего интерфейсный тип storage.StorageMock
	strg, err := storage.NewStorage(log.Logger, flags.DBStorageAdress, flags.FileStoragePath)
	if err != nil {
		return nil, err
	}
	// Create a handlers object.
	hndlrs := handlers.NewHandlers(strg, log.Logger, flags.Prefix)
	// Create a router object to pass to http.ListenAndServe.
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

// Delete ShortLink continuously listens to the a.Hndlrs.DelCn channel.
// Which receives bundles of shortened URLs for their deletion
func (a *app) DeleteShortLink() {
	ticker := time.NewTicker(1 * time.Second)

	var shortLinkList []*models.DeleteURL
	for {
		select {
		// The channel receives data in the form of a models structure.DeleteURL
		case shortLinksStruct := <-a.Hndlrs.DelCn:
			shortLinkList = append(shortLinkList, shortLinksStruct)
		// Each ticker function checks the received shortLinksStruct in the shortLinkList.
		// If the shortLinkList is not empty, the program compiles two lists of userID and shortLink from this data
		// and passes them to the UpdateDeletedFlag function.
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
