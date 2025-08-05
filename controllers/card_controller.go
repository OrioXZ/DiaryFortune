package controllers

import (
	"net/http"

	"dairyfortune/config"
	"dairyfortune/models"
	"dairyfortune/utils"

	"github.com/gin-gonic/gin"
)

func GetCards(c *gin.Context) {
	status := c.Query("status")
	name := c.Query("name")
	username := c.Query("username")

	if !utils.IsAdmin(username) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var cards []models.Card
	query := config.DB

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if err := query.Find(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cards"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":  len(cards),
		"result": cards,
	})
}

func UpdateCard(c *gin.Context) {
	var input CardUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if !utils.IsAdmin(input.Username) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var card models.Card
	if err := config.DB.Where("name = ?", input.Name).First(&card).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	card.Message = input.Message
	card.Type = input.Type
	card.Rarity = input.Rarity
	card.ImagePath = input.ImagePath
	card.Status = input.Status

	if err := config.DB.Save(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card updated", "card": card})
}

func CreateCard(c *gin.Context) {
	var input CardUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if !utils.IsAdmin(input.Username) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Check for duplicate name
	var existing models.Card
	if err := config.DB.Where("name = ?", input.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Card with that name already exists"})
		return
	}

	card := models.Card{
		Name:      input.Name,
		Message:   input.Message,
		Type:      input.Type,
		Rarity:    input.Rarity,
		ImagePath: input.ImagePath,
		Status:    input.Status,
	}

	if err := config.DB.Create(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create card"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Card created", "card": card})
}

func DeleteCard(c *gin.Context) {
	var input CardDeleteInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if !utils.IsAdmin(input.Username) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var card models.Card
	if err := config.DB.Where("name = ?", input.Name).First(&card).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	card.Status = "N"
	if err := config.DB.Save(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark card inactive"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card marked as inactive", "card": card})
}

type CardUpdateInput struct {
	Username  string `json:"username"`
	Name      string `json:"name" gorm:"not null"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	Rarity    string `json:"rarity"`
	ImagePath string `json:"imagePath"`
	Status    string `json:"status"`
}

type CardDeleteInput struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}
