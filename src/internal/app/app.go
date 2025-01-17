package app

import (
	"log"
	"root/internal/config"

	"golang.org/x/sync/errgroup"
)

type App struct {
	config          *config.Config
	httpServer      *httpServer
	websocketServer *websocketServer
}

func NewApp() (*App, error) {
	app := &App{
		config: config.GetConfig(),
	}

	storageProvider := NewStorageProvider(app.config.Database.URL)

	serviceProvider := NewServiceProvider(
		app.config,
		storageProvider,
	)

	app.httpServer = NewHttpServer(
		app.config,
		serviceProvider,
	)

	app.websocketServer = NewWebsocketServer(
		app.httpServer.Echo(),
		storageProvider.Storage(),
	)

	return app, nil
}

func (a *App) Run() error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		log.Println("Websocker server is running")
		return a.websocketServer.Run()
	})

	eg.Go(func() error {
		return a.httpServer.Run()
	})

	return eg.Wait()
}
