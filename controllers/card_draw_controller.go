package controllers

import (
	"encoding/json"
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
		var cachedCard models.Card
		if jsonErr := json.Unmarshal([]byte(val), &cachedCard); jsonErr == nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "You already drew a card today",
				"card":    cachedCard,
				"source":  "cache",
			})
			return
		}
		// If unmarshal fails, fallback to DB logic
	}

	// 🔁 2. Proceed to DB logic
	user := models.User{Username: username}
	if err := config.DB.Where("username = ?", username).FirstOrCreate(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find/create user"})
		return
	}

	// 3. Check if already drawn today
	var existingDraw models.CardDraw
	today := time.Now().Truncate(24 * time.Hour)
	err = config.DB.
		Preload("User").
		Preload("Card").
		Where("user_id = ? AND date = ?", user.ID, today).
		First(&existingDraw).Error

	if err == nil {
		// ✅ Cache the existing draw for next time
		jsonBytes, _ := json.Marshal(existingDraw.Card)
		_ = config.Rdb.Set(config.Ctx, "draw:"+username, string(jsonBytes), 24*time.Hour).Err()

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

	// ✅ Cache the new draw
	jsonBytes, _ := json.Marshal(card)
	_ = config.Rdb.Set(config.Ctx, "draw:"+username, string(jsonBytes), 24*time.Hour).Err()

	// 7. Return result
	c.JSON(http.StatusOK, gin.H{
		"message": "Card drawn!",
		"card":    card,
		"source":  "fresh-db",
	})
}
