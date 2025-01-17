package storage

import "root/internal/domain"

type UserStorage interface {
	GetByID(id int) (*domain.User, error)
	GetByLogin(login string) (*domain.User, error)
	Create(user *domain.User) error
	Delete(id int) error
}
