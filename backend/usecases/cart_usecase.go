package usecases

import (
	"context"
	"errors"

	"github.com/garvpatidar/food-delivery/backend/models"
	"github.com/garvpatidar/food-delivery/backend/repositories"
)

type CartUsecase interface {
	AddItem(ctx context.Context, userID uint, item models.CartItem) (*models.Cart, error)
	RemoveItem(ctx context.Context, userID uint, foodItemID uint) (*models.Cart, error)
	GetCart(ctx context.Context, userID uint) (*models.Cart, error)
	ClearCart(ctx context.Context, userID uint) error
}

type cartUsecase struct {
	cartRepo repositories.CartRepository
}

func NewCartUsecase(cartRepo repositories.CartRepository) CartUsecase {
	return &cartUsecase{cartRepo: cartRepo}
}

func (u *cartUsecase) AddItem(ctx context.Context, userID uint, item models.CartItem) (*models.Cart, error) {
	// Business Rule: Quantity must be valid
	if item.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	// 1. Pehle user ki existing cart uthao
	cart, err := u.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	cart.UserID = userID

	// 2. Check karo ki item pehle se cart mein hai ya nahi
	found := false
	for i, existingItem := range cart.Items {
		if existingItem.FoodItemID == item.FoodItemID {
			// Agar hai, toh sirf quantity badhao
			cart.Items[i].Quantity += item.Quantity
			found = true
			break
		}
	}

	// 3. Agar nahi tha, toh add karo
	if !found {
		cart.Items = append(cart.Items, item)
	}

	// 4. Total recalculate karo
	cart.TotalAmount = calculateTotal(cart.Items)

	// 5. Updated cart Redis mein save karo
	if err := u.cartRepo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (u *cartUsecase) RemoveItem(ctx context.Context, userID uint, foodItemID uint) (*models.Cart, error) {
	cart, err := u.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Filter out the item to remove
	updatedItems := []models.CartItem{}
	for _, item := range cart.Items {
		if item.FoodItemID != foodItemID {
			updatedItems = append(updatedItems, item)
		}
	}

	cart.Items = updatedItems
	cart.TotalAmount = calculateTotal(cart.Items)

	if err := u.cartRepo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (u *cartUsecase) GetCart(ctx context.Context, userID uint) (*models.Cart, error) {
	return u.cartRepo.GetCart(ctx, userID)
}

func (u *cartUsecase) ClearCart(ctx context.Context, userID uint) error {
	return u.cartRepo.ClearCart(ctx, userID)
}

// Helper: Total calculate karta hai
func calculateTotal(items []models.CartItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}
