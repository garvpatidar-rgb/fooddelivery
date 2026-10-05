package usecases

import (
	"context"
	"errors"

	"github.com/garvpatidar/food-delivery/backend/models"
	"github.com/garvpatidar/food-delivery/backend/repositories"
)

// RestaurantUsecase defines the interface for restaurant business logic
type RestaurantUsecase interface {
	CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) error
	GetRestaurant(ctx context.Context, id uint) (*models.Restaurant, error)
	GetAllRestaurants(ctx context.Context) ([]models.Restaurant, error)
}

type restaurantUsecase struct {
	restaurantRepo repositories.RestaurantRepository
}

// NewRestaurantUsecase creates a new instance of RestaurantUsecase
func NewRestaurantUsecase(repo repositories.RestaurantRepository) RestaurantUsecase {
	return &restaurantUsecase{
		restaurantRepo: repo,
	}
}

func (u *restaurantUsecase) CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) error {
	// Business Rule: Name cannot be empty
	if restaurant.Name == "" {
		return errors.New("restaurant name cannot be empty")
	}

	// Business Rule: Address cannot be empty
	if restaurant.Address == "" {
		return errors.New("restaurant address cannot be empty")
	}

	return u.restaurantRepo.Create(ctx, restaurant)
}

func (u *restaurantUsecase) GetRestaurant(ctx context.Context, id uint) (*models.Restaurant, error) {
	if id == 0 {
		return nil, errors.New("invalid restaurant ID")
	}
	return u.restaurantRepo.GetByID(ctx, id)
}

func (u *restaurantUsecase) GetAllRestaurants(ctx context.Context) ([]models.Restaurant, error) {
	return u.restaurantRepo.GetAll(ctx)
}
