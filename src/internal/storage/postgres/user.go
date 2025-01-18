package postgres

import (
	"root/internal/domain"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type userStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *userStorage {
	return &userStorage{db: db}
}

func (s *userStorage) GetByID(userID int) (*domain.User, error) {
	user := new(domain.User)
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}
	return user, nil
}

func (s *userStorage) GetByLogin(login string) (*domain.User, error) {
	user := new(domain.User)
	if err := s.db.First(&user, "login = ?", login).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}
	return user, nil
}

func (s *userStorage) Create(user *domain.User) error {
	if err := s.db.Create(user).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *userStorage) Delete(userID int) error {
	if err := s.db.Delete(&domain.User{Id: userID}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}
