package app

import (
	"log"
	"root/internal/storage"
	"root/internal/storage/postgres"
	db "root/pkg/database"
)

type storageProvider struct {
	storage storage.Storage
}

func NewStorageProvider(dbURL string) *storageProvider {
	db, err := db.NewConnect(dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return &storageProvider{
		storage: storage.NewStorage(
			postgres.NewUserStorage(db),
		),
	}
}

func (s *storageProvider) Storage() storage.Storage {
	return s.storage
}
