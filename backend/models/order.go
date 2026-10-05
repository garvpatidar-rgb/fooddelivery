package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"not null" json:"user_id"`
	User            User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	RestaurantID    uint           `gorm:"not null" json:"restaurant_id"` // NEW: Order belongs to a restaurant
	Restaurant      Restaurant     `gorm:"foreignKey:RestaurantID" json:"restaurant,omitempty"`
	Items           []OrderItem    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items"`
	TotalAmount     float64        `gorm:"not null" json:"total_amount"`
	Status          string         `gorm:"size:30;default:'pending'" json:"status"` // pending, preparing, out_for_delivery, delivered, cancelled
	DeliveryAddress string         `gorm:"type:text;not null" json:"delivery_address"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderItem struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	OrderID    uint     `gorm:"not null" json:"order_id"`
	FoodItemID uint     `gorm:"not null" json:"food_item_id"`
	FoodItem   FoodItem `gorm:"foreignKey:FoodItemID" json:"food_item,omitempty"`
	Quantity   int      `gorm:"not null" json:"quantity"`
	Price      float64  `gorm:"not null" json:"price"` // Price at time of order
}
