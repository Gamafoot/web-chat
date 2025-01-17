package v1

import (
	"errors"
	"net/http"
	"root/internal/domain"
	"root/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

func (h *handler) LoginPage(c echo.Context) error {
	data := make(map[string]string, 0)

	if c.Request().Method == http.MethodPost {
		resp := h.loginOrCreate(c)
		if resp.Err != nil {
			return resp.Err
		}

		if resp.Success {
			return c.Redirect(http.StatusFound, "/")
		} else {
			data["Error"] = resp.Message
		}
	}

	return c.Render(http.StatusOK, "login.html", data)
}

type (
	credentialsInput struct {
		Login    string `json:"login" validate:"required,gte=1,lte=15"`
		Password string `json:"password" validate:"required"`
	}

	loginOrCreateResp struct {
		Success bool
		Message string
		Err     error
	}
)

func (h *handler) loginOrCreate(c echo.Context) loginOrCreateResp {
	input := &credentialsInput{
		Login:    c.FormValue("login"),
		Password: c.FormValue("password"),
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return loginOrCreateResp{
			Message: "Validation failed, min length of username is 1, max is 15",
		}
	}

	err := h.authService.LoginOrCreate(c, service.LoginInput{
		Login:    input.Login,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidLoginOrPassword) {
			return loginOrCreateResp{
				Message: "Invalid login or password",
			}
		}

		return loginOrCreateResp{
			Err: err,
		}
	}

	return loginOrCreateResp{
		Success: true,
	}
}

func (h *handler) GetSession(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	userId := sess.Values["userId"].(int)

	return c.JSON(http.StatusOK, map[string]int{
		"userId": userId,
	})
}
