package service

import (
	"root/internal/config"
	"root/internal/domain"
	"root/internal/storage"
	"root/pkg/hash"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

type AuthService struct {
	config  *config.Config
	storage storage.Storage
}

func NewAuthService(config *config.Config, storage storage.Storage) *AuthService {
	return &AuthService{
		config:  config,
		storage: storage,
	}
}

type (
	LoginInput struct {
		Login    string
		Password string
	}
)

func (s *AuthService) LoginOrCreate(c echo.Context, input LoginInput) error {
	var (
		user *domain.User
	)

	hasher := hash.NewSHA1Hasher(s.config.Hash.Salt)

	user, err := s.storage.User.GetByLogin(input.Login)
	if err != nil {
		hash, err := hasher.Hash(input.Password)
		if err != nil {
			return pkgErrors.WithStack(err)
		}

		err = s.storage.User.Create(&domain.User{
			Login:    input.Login,
			Password: hash,
		})
		if err != nil {
			return pkgErrors.WithStack(err)
		}

		user, err = s.storage.User.GetByLogin(input.Login)
		if err != nil {
			return pkgErrors.WithStack(err)
		}
	} else {
		hash, err := hasher.Hash(input.Password)
		if err != nil {
			return pkgErrors.WithStack(err)
		}

		if user.Password != hash {
			return domain.ErrInvalidLoginOrPassword
		}
	}

	return s.createSession(c, user.Id)
}

func (s *AuthService) createSession(c echo.Context, userId int) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	sess.Values["userId"] = userId
	sess.Save(c.Request(), c.Response())

	return nil
}
