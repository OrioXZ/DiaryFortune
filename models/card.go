package models

type Card struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"not null"`      // Name of the card (e.g., The Sun, Death)
	Message   string `json:"message" gorm:"type:text"`  // Fortune text
	Type      string `json:"type" gorm:"not null"`      // good / bad / secret
	Rarity    string `json:"rarity"`                    // common / rare / legendary
	ImagePath string `json:"imagePath"`                 // optional image path
	Status    string `json:"status" gorm:"default:'Y'"` // Y = active, N = inactive
}
