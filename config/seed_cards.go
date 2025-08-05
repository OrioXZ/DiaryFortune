package config

import "dairyfortune/models"

func SeedCards() {
	var count int64
	DB.Model(&models.Card{}).Count(&count)
	if count > 0 {
		return // Cards already exist, skip seeding
	}

	cards := []models.Card{
		{
			Name:    "The Sun",
			Message: "Success and joy shine upon you.",
			Type:    "good",
			Rarity:  "common",
			Status:  "Y",
		},
		{
			Name:    "The Fool",
			Message: "Take a leap of faith — the universe has your back.",
			Type:    "good",
			Rarity:  "common",
			Status:  "Y",
		},
		{
			Name:    "The Tower",
			Message: "Disruption leads to transformation. Be strong.",
			Type:    "bad",
			Rarity:  "common",
			Status:  "Y",
		},
		{
			Name:    "The Moon",
			Message: "Confusion is temporary. Trust your intuition.",
			Type:    "bad",
			Rarity:  "common",
			Status:  "Y",
		},
		{
			Name:    "The Star",
			Message: "Hope will guide your way.",
			Type:    "good",
			Rarity:  "rare",
			Status:  "Y",
		},
		{
			Name:    "Shadow Whisper",
			Message: "A secret path has opened.",
			Type:    "secret",
			Rarity:  "legendary",
			Status:  "Y",
		},
	}

	DB.Create(&cards)
}
