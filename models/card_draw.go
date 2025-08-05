package models

import "time"

type CardDraw struct {
	ID     uint      `gorm:"primaryKey"`
	UserID uint      `json:"userId"`
	CardID uint      `json:"cardId"`
	Date   time.Time `json:"date"`

	User User `gorm:"foreignKey:UserID"`
	Card Card `gorm:"foreignKey:CardID"`
}
