package controllers

import (
	"math/rand"
	"net/http"
	"time"

	"dairyfortune/config"
	"dairyfortune/models"

	"github.com/gin-gonic/gin"
)

func DrawCard(c *gin.Context) {
	// 1. Get username from query or body
	username := c.Query("username")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	// 2. Find or create users
	user := models.User{Username: username}
	if err := config.DB.Where("username = ?", username).FirstOrCreate(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find/create user"})
		return
	}

	// 3. Check if already drawn today
	var existingDraw models.CardDraw
	today := time.Now().Truncate(24 * time.Hour)
	err := config.DB.
		Preload("User").
		Preload("Card").
		Where("user_id = ? AND date = ?", user.ID, today).
		First(&existingDraw).Error

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "You already drew a card today",
			"result":  existingDraw,
		})
		return
	}

	// 4. Get all cards
	var cards []models.Card
	if err := config.DB.Find(&cards).Error; err != nil || len(cards) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No cards available"})
		return
	}

	// 5. Randomly pick one
	rand.Seed(time.Now().UnixNano())
	card := cards[rand.Intn(len(cards))]

	// 6. Save draw
	draw := models.CardDraw{
		UserID: user.ID,
		CardID: card.ID,
		Date:   today,
	}
	if err := config.DB.Create(&draw).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save draw"})
		return
	}

	// 7. Return result
	c.JSON(http.StatusOK, gin.H{
		"message": "Card drawn!",
		"card":    card,
	})
}
