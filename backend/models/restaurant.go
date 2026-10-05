package models

import (
	"time"

	"gorm.io/gorm"
)

type Restaurant struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Address     string         `gorm:"type:text;not null" json:"address"`
	Phone       string         `gorm:"size:20" json:"phone"`
	OwnerID     uint           `gorm:"not null" json:"owner_id"` // Links to the User who owns it
	Owner       User           `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	FoodItems   []FoodItem     `gorm:"foreignKey:RestaurantID;constraint:OnDelete:CASCADE" json:"food_items,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
