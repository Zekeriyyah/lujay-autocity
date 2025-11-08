package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(databaseURL string) *gorm.DB {
	var dsn string

	// either using postgres on render or from local
	if os.Getenv("RENDER") == "true" {
		dsn = os.Getenv("DATABASE_URL")
	} else {
		dsn = databaseURL
	}	

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}

	DB = db
	log.Println("✅ Database connected successfully!")

	//Auto-migrate the db
	Run()
	return DB
}
