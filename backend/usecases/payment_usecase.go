package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/garvpatidar/food-delivery/backend/models"
	"github.com/garvpatidar/food-delivery/backend/repositories"
)

type PaymentUsecase interface {
	InitiatePayment(ctx context.Context, orderID uint, userID uint, amount float64, method string) (*models.Payment, error)
	ConfirmPayment(ctx context.Context, paymentID uint, transactionID string) error
	GetPaymentByOrder(ctx context.Context, orderID uint) (*models.Payment, error)
}

type paymentUsecase struct {
	paymentRepo repositories.PaymentRepository
}

func NewPaymentUsecase(repo repositories.PaymentRepository) PaymentUsecase {
	return &paymentUsecase{paymentRepo: repo}
}

func (u *paymentUsecase) InitiatePayment(ctx context.Context, orderID uint, userID uint, amount float64, method string) (*models.Payment, error) {
	// Business Rule: Amount must be positive
	if amount <= 0 {
		return nil, errors.New("payment amount must be greater than 0")
	}

	// Business Rule: Valid payment methods only
	validMethods := map[string]bool{"card": true, "upi": true, "netbanking": true, "cod": true}
	if !validMethods[method] {
		return nil, errors.New("invalid payment method")
	}

	payment := &models.Payment{
		OrderID:       orderID,
		UserID:        userID,
		Amount:        amount,
		Currency:      "INR",
		Status:        models.PaymentStatusPending,
		PaymentMethod: method,
		// In real Stripe integration: call stripe.PaymentIntents.Create() here
		// For now, we generate a mock transaction ID
		TransactionID: fmt.Sprintf("TXN_%d_%d", orderID, time.Now().Unix()),
	}

	if err := u.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (u *paymentUsecase) ConfirmPayment(ctx context.Context, paymentID uint, transactionID string) error {
	// In real Stripe: verify webhook signature and confirm payment intent
	// Here we simulate a successful payment confirmation
	return u.paymentRepo.UpdateStatus(ctx, paymentID, models.PaymentStatusSuccess, transactionID)
}

func (u *paymentUsecase) GetPaymentByOrder(ctx context.Context, orderID uint) (*models.Payment, error) {
	return u.paymentRepo.GetByOrderID(ctx, orderID)
}
