package models

import (
	"time"

	"gorm.io/gorm"
)

type FoodItem struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	RestaurantID uint           `gorm:"not null" json:"restaurant_id"` // NEW: Links food to a restaurant
	Name         string         `gorm:"size:100;not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Price        float64        `gorm:"not null" json:"price"`
	Category     string         `gorm:"size:50;not null" json:"category"` // e.g., Pizza, Burger, Drinks
	ImageURL     string         `gorm:"size:255" json:"image_url"`
	IsAvailable  bool           `gorm:"default:true" json:"is_available"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
