package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/garvpatidar/food-delivery/backend/models"
	"github.com/redis/go-redis/v9"
)

type CartRepository interface {
	SaveCart(ctx context.Context, cart *models.Cart) error
	GetCart(ctx context.Context, userID uint) (*models.Cart, error)
	ClearCart(ctx context.Context, userID uint) error
}

type cartRepository struct {
	redisClient *redis.Client
	memoryStore map[uint]*models.Cart
	mu          sync.RWMutex
}

func NewCartRepository(redisClient *redis.Client) CartRepository {
	return &cartRepository{
		redisClient: redisClient,
		memoryStore: make(map[uint]*models.Cart),
	}
}

func (r *cartRepository) SaveCart(ctx context.Context, cart *models.Cart) error {
	if r.redisClient == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.memoryStore[cart.UserID] = cart
		return nil
	}

	key := fmt.Sprintf("cart:%d", cart.UserID)
	cartJSON, err := json.Marshal(cart)
	if err != nil {
		return err
	}
	return r.redisClient.Set(ctx, key, cartJSON, 48*time.Hour).Err()
}

func (r *cartRepository) GetCart(ctx context.Context, userID uint) (*models.Cart, error) {
	if r.redisClient == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		if cart, exists := r.memoryStore[userID]; exists {
			return cart, nil
		}
		return &models.Cart{UserID: userID, Items: []models.CartItem{}, TotalAmount: 0}, nil
	}

	key := fmt.Sprintf("cart:%d", userID)
	cartJSON, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return &models.Cart{UserID: userID, Items: []models.CartItem{}, TotalAmount: 0}, nil
	} else if err != nil {
		return nil, err
	}

	var cart models.Cart
	if err := json.Unmarshal([]byte(cartJSON), &cart); err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) ClearCart(ctx context.Context, userID uint) error {
	if r.redisClient == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		delete(r.memoryStore, userID)
		return nil
	}

	key := fmt.Sprintf("cart:%d", userID)
	return r.redisClient.Del(ctx, key).Err()
}
