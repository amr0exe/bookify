package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dbURL string) (*gorm.DB, error) {
	db, err := InitDB(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed db connection: %w", err)
	}

	if err := CheckDBHealth(db); err != nil {
		return nil, fmt.Errorf("healthcheck failed: %w", err)
	}

	return db, nil
}

func InitDB(db_url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(db_url), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	psDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := psDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db health check failed (ping timeout or connection refused): %w", err)
	}

	log.Println("DB connection success & healthy !!!")
	return db, nil
}

func CheckDBHealth(db *gorm.DB) error {
	var result int

	err := db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("unexpected healthcheck result")
	}

	return nil
}
