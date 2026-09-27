package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kjwardy/acronom/api/app"
	"github.com/kjwardy/acronom/api/config"
	"github.com/kjwardy/acronom/api/db"
	"github.com/labstack/echo/v4"
)

func main() {
	appConfig, err := config.Load()
	if err != nil {
		log.Fatalf("invalid application configuration: %v", err)
	}

	connectionContext, cancelConnection := context.WithTimeout(context.Background(), 10*time.Second)
	database, err := db.Connect(connectionContext, appConfig.Database)
	cancelConnection()
	if err != nil {
		log.Fatalf("database startup check failed: %v", err)
	}
	defer database.Close()

	e := echo.New()
	e.Debug = appConfig.Debug
	app.Init(e, database)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", appConfig.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Error(err)
	}
}
