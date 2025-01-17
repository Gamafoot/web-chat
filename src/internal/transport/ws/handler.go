package ws

import (
	"root/internal/storage"
	"root/internal/transport/middleware"

	"github.com/labstack/echo/v4"
)

type handler struct {
	hub     *Hub
	storage storage.Storage
}

func NewHandler(hub *Hub, storage storage.Storage) *handler {
	return &handler{
		hub:     hub,
		storage: storage,
	}
}

func (h *handler) InitRoutes(e *echo.Echo) {
	g := e.Group("/ws", middleware.RequiredAuth)
	{
		g.GET("/joinRoom/:name", h.JoinRoom)
		g.GET("/getClients/:name", h.GetClients)
	}
}
