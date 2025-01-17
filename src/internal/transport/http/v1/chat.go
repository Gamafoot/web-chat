package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) ChatPage(c echo.Context) error {
	return c.Render(http.StatusOK, "chat.html", nil)
}

func (h *handler) PingPong(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
