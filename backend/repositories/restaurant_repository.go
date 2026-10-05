package repositories

import (
	"context"

	"github.com/garvpatidar/food-delivery/backend/models"
	"gorm.io/gorm"
)

// RestaurantRepository defines the interface for restaurant database operations
type RestaurantRepository interface {
	Create(ctx context.Context, restaurant *models.Restaurant) error
	GetByID(ctx context.Context, id uint) (*models.Restaurant, error)
	GetAll(ctx context.Context) ([]models.Restaurant, error)
	Update(ctx context.Context, restaurant *models.Restaurant) error
	Delete(ctx context.Context, id uint) error
}

type restaurantRepository struct {
	db *gorm.DB
}

// NewRestaurantRepository creates a new instance of RestaurantRepository
func NewRestaurantRepository(db *gorm.DB) RestaurantRepository {
	return &restaurantRepository{db: db}
}

func (r *restaurantRepository) Create(ctx context.Context, restaurant *models.Restaurant) error {
	return r.db.WithContext(ctx).Create(restaurant).Error
}

func (r *restaurantRepository) GetByID(ctx context.Context, id uint) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	err := r.db.WithContext(ctx).First(&restaurant, id).Error
	return &restaurant, err
}

func (r *restaurantRepository) GetAll(ctx context.Context) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := r.db.WithContext(ctx).Find(&restaurants).Error
	return restaurants, err
}

func (r *restaurantRepository) Update(ctx context.Context, restaurant *models.Restaurant) error {
	return r.db.WithContext(ctx).Save(restaurant).Error
}

func (r *restaurantRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Restaurant{}, id).Error
}
