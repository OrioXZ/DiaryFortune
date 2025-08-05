package controllers

import (
	"dairyfortune/config"
	"dairyfortune/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// This is a dummy router to test draw handler
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/draw", DrawCard)
	return r
}

func TestDrawCard(t *testing.T) {
	config.InitTestDB()
	config.InitRedis()

	// 1. Create test user
	user := models.User{Username: "testuser"}
	config.DB.Create(&user)

	// 2. Create test card
	card := models.Card{Name: "Test Card"}
	config.DB.Create(&card)

	// 3. Create test card draw for TODAY
	draw := models.CardDraw{
		UserID: user.ID,
		CardID: card.ID,
		Date:   time.Now().Truncate(24 * time.Hour), // match your controller logic
	}
	config.DB.Create(&draw)
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/draw?username=testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "card")
}
