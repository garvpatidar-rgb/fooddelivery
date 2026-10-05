package database

import (
	"log"

	"github.com/garvpatidar/food-delivery/backend/config"
	"github.com/garvpatidar/food-delivery/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB instance holds the connection to SQLite
var DB *gorm.DB

func ConnectDB(cfg config.Config) {
	// Connect to SQLite database file named "food_delivery.db"
	db, err := gorm.Open(sqlite.Open("food_delivery.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("SQLite Database connection established successfully!")
	DB = db

	// Auto Migrate Models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Restaurant{},
		&models.FoodItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
	)
	if err != nil {
		log.Fatalf("Failed to run auto migrations: %v", err)
	}

	log.Println("Database migration completed successfully!")

	// Seed initial data (Admin, Restaurants, Food Items)
	SeedData(DB)
}
