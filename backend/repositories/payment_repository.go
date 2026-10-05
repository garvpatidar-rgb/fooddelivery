package repositories

import (
	"context"

	"github.com/garvpatidar/food-delivery/backend/models"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) error
	GetByOrderID(ctx context.Context, orderID uint) (*models.Payment, error)
	UpdateStatus(ctx context.Context, id uint, status string, transactionID string) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id uint, status string, transactionID string) error {
	return r.db.WithContext(ctx).Model(&models.Payment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         status,
			"transaction_id": transactionID,
		}).Error
}
