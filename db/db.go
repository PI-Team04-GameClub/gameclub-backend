package db

import (
	"log"

	"github.com/PI-Team04-GameClub/gameclub-backend/config"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	db, err := gorm.Open(postgres.Open(cfg.ConnString()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = db
	log.Println("Database connected successfully")
}

func Migrate() {
	if err := DB.AutoMigrate(
		&models.Game{},
		&models.User{},
		&models.Team{},
		&models.Tournament{},
		&models.News{},
		&models.Comment{},
		&models.FriendRequest{},
	); err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}
	log.Println("Database migration completed successfully")
}
