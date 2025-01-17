package main

import (
	"log"
	"root/internal/config"
	"root/internal/domain"
	db "root/pkg/database"
)

func main() {
	log.Println("start migrating...")

	cfg := config.GetConfig()

	db, err := db.NewConnect(cfg.Database.URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	tables := []interface{}{
		domain.User{},
	}

	if err = db.AutoMigrate(tables...); err != nil {
		log.Fatalf("failed to migrate models: %v", err)
	}

	log.Println("migrate was successed")
}
