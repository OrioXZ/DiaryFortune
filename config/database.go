package config

import (
	"dairyfortune/models"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("fortunes.db"), &gorm.Config{})
	if err != nil {
		panic("❌ Failed to connect to database: " + err.Error())
	}

	fmt.Println("✅ Connected to database")
	DB = database
}

func InitTestDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test SQLite DB")
	}

	// Auto-migrate your models here
	DB.AutoMigrate(
		&models.User{},
		&models.Card{},
		&models.CardDraw{},
		&models.Achievement{},
	)

}
