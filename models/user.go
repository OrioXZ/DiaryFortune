package models

type User struct {
	ID       uint   `gorm:"primaryKey"`                      // Auto-increment by default
	Username string `json:"username" gorm:"unique;not null"` // Must be unique
	Status   string `json:"status" gorm:"default:'Y'"`       // Default is "Y" if not provided
}
