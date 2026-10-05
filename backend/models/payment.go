package models

import (
	"time"

	"gorm.io/gorm"
)

// Payment statuses
const (
	PaymentStatusPending  = "pending"
	PaymentStatusSuccess  = "success"
	PaymentStatusFailed   = "failed"
	PaymentStatusRefunded = "refunded"
)

type Payment struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	OrderID         uint           `gorm:"not null;uniqueIndex" json:"order_id"` // One payment per order
	UserID          uint           `gorm:"not null" json:"user_id"`
	Amount          float64        `gorm:"not null" json:"amount"`
	Currency        string         `gorm:"size:10;default:'INR'" json:"currency"`
	Status          string         `gorm:"size:30;default:'pending'" json:"status"` // pending, success, failed, refunded
	PaymentMethod   string         `gorm:"size:50" json:"payment_method"`           // card, upi, netbanking
	TransactionID   string         `gorm:"size:255;uniqueIndex" json:"transaction_id"` // From payment gateway
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
