package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kjwardy/acronom/api/app"
	"github.com/kjwardy/acronom/api/config"
	"github.com/labstack/echo/v4"
)

func main() {
	appConfig, err := config.Load()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Debug = appConfig.Debug
	app.Init(e)

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
