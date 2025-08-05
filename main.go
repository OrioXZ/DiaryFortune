package main

import (
	"dairyfortune/config"
	"dairyfortune/models"
	"dairyfortune/routes"
)

func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(
		&models.Card{},
		&models.User{},
		&models.CardDraw{},
		&models.Achievement{},
	)
	config.SeedCards()

	r := routes.SetupRouter()

	r.Run(":8080") // Starts server on localhost:8080
}
