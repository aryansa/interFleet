package src

import (
	"context"
	"errors"
	"fmt"
	"interFleet/src/device/handlers"
	"interFleet/src/device/infrastructure"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	HostURLConfig string
	CSVPathConfig string
}

func (app *App) Run(ctx context.Context) (err error) {
	if app.HostURLConfig == "" {
		app.HostURLConfig = "0.0.0.0:6733"
	}
	if app.CSVPathConfig == "" {
		app.CSVPathConfig = "./devices.csv"
	}

	repository, repoErr := infrastructure.NewDeviceRepository(app.CSVPathConfig)
	if repoErr != nil {
		return repoErr
	}

	handler := handlers.CommandHandler{Repo: repository}

	controller := infrastructure.Controller{Repo: repository, Handler: handler}

	gin.SetMode(gin.ReleaseMode)
	ginHandler := gin.Default()

	apiRoutes := ginHandler.Group("/api/v1")

	controller.RegisterRoutes(func(method string, url string, handler gin.HandlerFunc) {
		apiRoutes.Handle(method, url, handler)
	})

	server := &http.Server{
		Addr:    app.HostURLConfig,
		Handler: ginHandler,
	}

	context.AfterFunc(ctx, func() {
		fmt.Println("graceful shutdown...")
		if err == nil {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
				log.Println("shutdown error:", shutdownErr)
			}
		}
	})

	fmt.Printf("Listening on %s\n", app.HostURLConfig)
	if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil

}
