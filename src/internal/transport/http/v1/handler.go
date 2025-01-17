package v1

import (
	"root/internal/config"
	"root/internal/service"
	"root/internal/transport/middleware"

	"github.com/labstack/echo/v4"
)

type handler struct {
	config      *config.Config
	authService *service.AuthService
}

func NewHandler(config *config.Config, userService *service.AuthService) *handler {
	return &handler{
		config:      config,
		authService: userService,
	}
}

func (h *handler) InitRoutes(e *echo.Echo) {
	requiredAuth := e.Group("", middleware.RequiredAuth)
	{
		requiredAuth.GET("/", h.ChatPage)
		requiredAuth.GET("/get-session", h.GetSession)
		requiredAuth.POST("/ping", h.PingPong)
	}

	e.GET("/login", h.LoginPage)
	e.POST("/login", h.LoginPage)
}
