package database

import (
	"log"

	"github.com/zekeriyyah/lujay-autocity/internal/models"
)


func Run() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Vehicle{},
		&models.Inspection{},
		&models.Listing{},
		&models.Transaction{},
		&models.Image{},
	)

	if err != nil {
		log.Fatal("❌ database auto-migration failed:", err)
	}
	log.Println("✅ Migrations completed successfully!")
}