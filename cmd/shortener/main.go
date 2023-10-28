// Package main is the entry point of the application.
// It initializes the application and starts the server.
package main

import (
	"go.uber.org/zap"
	"net/http"
	"practicumserver/internal/app"
	"practicumserver/internal/storage/pg"
	"practicumserver/internal/utils"
)

func main() {
	// Beginning of the program execution.
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	// Create a new application instance.
	appStrct, err := app.NewApp()
	if err != nil {
		return err
	}
	// Close log files when done.
	defer utils.Closelog(appStrct.Log, appStrct.Flags)

	// Close the database connection if it's a database storage.
	defer func() {
		if DBstrg, ok := appStrct.Strg.(*pgstorage.DBStorage); ok {
			DBstrg.DB.Close()
		}
	}()

	// Start the server at the address specified in flags.String().
	appStrct.Log.Info("Server start", zap.String("Running server on:", appStrct.Flags.String()))
	err = http.ListenAndServe(appStrct.Flags.String(), appStrct.Rtr)
	appStrct.Log.Error("Error:", zap.Error(err))
	return err
}
