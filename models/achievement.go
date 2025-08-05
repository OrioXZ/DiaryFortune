package models

import "time"

type Achievement struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `json:"userId"`
	Name       string    `json:"name"` // e.g., "Good Luck Strike"
	AchievedAt time.Time `json:"achievedAt"`

	User User `gorm:"foreignKey:UserID"`
}
