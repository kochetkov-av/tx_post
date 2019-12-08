package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/shopspring/decimal"
)

func initDb(config *Config) (*gorm.DB, error) {
	// DB initialization.
	db, err := gorm.Open("postgres", config.GetConnectionString())
	if err != nil {
		return nil, err

	}

	db.AutoMigrate(&Transaction{}, &User{})

	// User initialization.
	var count int
	db.Model(&User{}).Count(&count)
	// Create user if not exist.
	if count < 1 {
		db.Create(&User{Balance: decimal.NewFromInt(0)})
	}

	return db, nil
}
