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
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}

	// 🔍 1. Try Redis cache first
	val, err := config.Rdb.Get(config.Ctx, "draw:"+username).Result()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"card":   val,
			"source": "cache",
		})
		return
	}

	// 🔁 2. Proceed to DB logic
	user := models.User{Username: username}
	if err := config.DB.Where("username = ?", username).FirstOrCreate(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find/create user"})
		return
	}

	// 3. Check if already drawn today
	today := time.Now().Truncate(24 * time.Hour)
	var existingDraw models.CardDraw
	err = config.DB.
		Preload("User").
		Preload("Card").
		Where("user_id = ? AND date = ?", user.ID, today).
		First(&existingDraw).Error

	if err == nil {
		// ✅ Cache it now for next time
		cacheValue := existingDraw.Card.Name // or marshal full card as JSON
		_ = config.Rdb.Set(config.Ctx, "draw:"+username, cacheValue, 0).Err()

		c.JSON(http.StatusOK, gin.H{
			"message": "You already drew a card today",
			"result":  existingDraw,
			"source":  "fresh-db",
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

	// ✅ Cache the drawn card for next time
	cacheValue := card.Name // or marshal as JSON for full card
	_ = config.Rdb.Set(config.Ctx, "draw:"+username, cacheValue, 0).Err()

	// 7. Return result
	c.JSON(http.StatusOK, gin.H{
		"message": "Card drawn!",
		"card":    card,
		"source":  "fresh-db",
	})
}
