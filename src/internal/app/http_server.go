package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"root/internal/config"
	v1 "root/internal/transport/http/v1"
	customMiddleware "root/internal/transport/middleware"
	"text/template"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type httpServer struct {
	config          *config.Config
	echo            *echo.Echo
	serviceProvider *serviceProvider
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewHttpServer(config *config.Config, serviceProvider *serviceProvider) *httpServer {
	server := &httpServer{
		config:          config,
		echo:            echo.New(),
		serviceProvider: serviceProvider,
	}

	server.echo.Static("/static", "assets/static")

	template := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("assets/html/*.html")),
	}
	server.echo.Renderer = template

	server.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: server.config.CORS.Origins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))
	server.echo.Use(customMiddleware.Logger)
	server.echo.Use(session.Middleware(sessions.NewCookieStore([]byte(config.Auth.SecretKey))))
	server.echo.Use(customMiddleware.BlockPathsIfAuth("/login"))

	handlerV1 := v1.NewHandler(
		server.config,
		server.serviceProvider.UserService(),
	)
	handlerV1.InitRoutes(server.echo)

	return server
}

func (s *httpServer) Run() error {
	addr := fmt.Sprintf(":%s", s.config.Listen.Port)

	if err := s.echo.Start(addr); err != nil {
		return err
	}

	return nil
}

func (s *httpServer) ShutdownServer(ctx context.Context) error {
	if err := s.echo.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (s *httpServer) Echo() *echo.Echo {
	return s.echo
}
