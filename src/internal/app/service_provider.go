package app

import (
	"root/internal/config"
	"root/internal/service"
)

type serviceProvider struct {
	config          *config.Config
	storageProvider *storageProvider

	userService *service.AuthService
}

func NewServiceProvider(config *config.Config, sp *storageProvider) *serviceProvider {
	return &serviceProvider{
		config:          config,
		storageProvider: sp,
	}
}

func (s *serviceProvider) UserService() *service.AuthService {
	if s.userService == nil {
		s.userService = service.NewAuthService(s.config, s.storageProvider.Storage())
	}

	return s.userService
}
