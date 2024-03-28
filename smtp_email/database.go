package main

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type EmailLog struct {
	gorm.Model
	From      string
	To        string
	Timestamp time.Time
	Success   bool
	Error     string
}

func initializeDB() *gorm.DB {
	// Read configuration from config.yaml
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// Read database configuration
	dbConfig := viper.GetStringMapString("database")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbConfig["host"],
		dbConfig["user"],
		dbConfig["password"],
		dbConfig["dbname"],
		dbConfig["port"],
		dbConfig["sslmode"],
		dbConfig["timezone"],
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto Migrate the Email model
	db.AutoMigrate(&EmailLog{})

	return db
}
