package middleware

import (
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			log.Printf("%+v", err)

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Internal Server Error",
			})
		}

		log.Printf("%s %s %d", c.Request().Method, c.Request().RequestURI, c.Response().Status)

		return nil
	}
}

func RequiredAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return pkgErrors.WithStack(err)
		}

		userId, ok := sess.Values["userId"].(int)
		if !ok {
			return c.Redirect(http.StatusFound, "/login")
		}

		c.Set("userId", userId)

		return next(c)
	}
}

func BlockPathsIfAuth(paths ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("session", c)
			if err != nil {
				return pkgErrors.WithStack(err)
			}

			_, ok := sess.Values["userId"].(int)
			if ok {
				for _, path := range paths {
					if c.Request().URL.Path == path {
						return c.Redirect(http.StatusFound, "/")
					}
				}
			}

			return next(c)
		}
	}
}
