package app

import (
	"root/internal/storage"
	"root/internal/transport/ws"

	"github.com/labstack/echo/v4"
)

type websocketServer struct {
	echo *echo.Echo
	hub  *ws.Hub
}

func NewWebsocketServer(echo *echo.Echo, storage storage.Storage) *websocketServer {
	s := &websocketServer{
		echo: echo,
		hub:  ws.NewHub(),
	}

	handler := ws.NewHandler(s.hub, storage)
	handler.InitRoutes(s.echo)

	return s
}

func (s *websocketServer) Run() error {
	s.hub.Run()
	return nil
}
