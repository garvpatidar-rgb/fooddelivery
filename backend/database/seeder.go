package database

import (
	"log"

	"github.com/garvpatidar/food-delivery/backend/models"
	"github.com/garvpatidar/food-delivery/backend/utils"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	// 1. Seed Admin & Customer Users if empty
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)

	var adminUser models.User
	var customerUser models.User

	if userCount == 0 {
		log.Println("Seeding initial users...")

		adminPass, _ := utils.HashPassword("admin123")
		adminUser = models.User{
			Name:     "Admin Manager",
			Email:    "admin@food.com",
			Password: adminPass,
			Role:     "admin",
			Phone:    "+91 9876543210",
			Address:  "TastyBites HQ, Tech Park",
		}
		db.Create(&adminUser)

		custPass, _ := utils.HashPassword("customer123")
		customerUser = models.User{
			Name:     "John Doe",
			Email:    "customer@food.com",
			Password: custPass,
			Role:     "customer",
			Phone:    "+91 9123456789",
			Address:  "Flat 402, Green Valley Apartments",
		}
		db.Create(&customerUser)

		log.Println("Default users seeded! Admin: admin@food.com | Customer: customer@food.com")
	} else {
		db.Where("role = ?", "admin").First(&adminUser)
	}

	// 2. Seed Default Restaurant if empty
	var restaurantCount int64
	db.Model(&models.Restaurant{}).Count(&restaurantCount)

	var defaultRestaurant models.Restaurant
	if restaurantCount == 0 {
		log.Println("Seeding initial restaurant...")

		ownerID := adminUser.ID
		if ownerID == 0 {
			ownerID = 1
		}

		defaultRestaurant = models.Restaurant{
			Name:        "TastyBites Grand Kitchen",
			Description: "Authentic gourmet burgers, wood-fired pizzas, fresh sushi & hand-crafted drinks.",
			Address:     "Central Food Street, Block B",
			Phone:       "+91 9988776655",
			OwnerID:     ownerID,
			IsActive:    true,
		}
		db.Create(&defaultRestaurant)
		log.Println("Default restaurant seeded!")
	} else {
		db.First(&defaultRestaurant)
	}

	// 3. Seed Food Items if empty
	var foodCount int64
	db.Model(&models.FoodItem{}).Count(&foodCount)

	if foodCount == 0 {
		log.Println("Seeding sample food items...")

		restID := defaultRestaurant.ID
		if restID == 0 {
			restID = 1
		}

		items := []models.FoodItem{
			{
				RestaurantID: restID,
				Name:         "Cheesy Deluxe Burger",
				Description:  "Juicy grilled patty topped with melted cheddar, crisp lettuce, secret sauce & toasted brioche bun.",
				Price:        249.00,
				Category:     "Burger",
				ImageURL:     "https://images.unsplash.com/photo-1568901346375-23c9450c58cd?w=600&auto=format&fit=crop&q=80",
				IsAvailable:  true,
			},
			{
				RestaurantID: restID,
				Name:         "Margherita Special Pizza",
				Description:  "Classic Naples sourdough crust with San Marzano tomato sauce, fresh mozzarella & sweet basil.",
				Price:        399.00,
				Category:     "Pizza",
				ImageURL:     "https://images.unsplash.com/photo-1604382354936-07c5d9983bd3?w=600&auto=format&fit=crop&q=80",
				IsAvailable:  true,
			},
			{
				RestaurantID: restID,
				Name:         "Creamy Alfredo Pasta",
				Description:  "Penne pasta tossed in rich parmesan garlic cream sauce with fresh herbs & garlic bread.",
				Price:        299.00,
				Category:     "Pasta",
				ImageURL:     "https://images.unsplash.com/photo-1621996346565-e3d5d6281290?w=600&auto=format&fit=crop&q=80",
				IsAvailable:  true,
			},
			{
				RestaurantID: restID,
				Name:         "Spicy Chicken Tacos",
				Description:  "3 crunchy corn tortillas filled with spiced shredded chicken, avocado salsa & chipotle crema.",
				Price:        279.00,
				Category:     "Mexican",
				ImageURL:     "https://images.unsplash.com/photo-1565299585323-38d6b0865b47?w=600&auto=format&fit=crop&q=80",
				IsAvailable:  true,
			},
			{
				RestaurantID: restID,
				Name:         "Dragon Sushi Roll",
				Description:  "Crispy tempura prawn roll topped with sliced avocado, spicy mayo and unagi reduction.",
				Price:        499.00,
				Category:     "Asian",
				ImageURL:     "https://images.unsplash.com/photo-1579871494447-9811cf80d66c?w=600&auto=format&fit=crop&q=80",
				IsAvailable:  true,
			},
			{
				RestaurantID: restID,
				Name:         "Chocolate Lava Cake",
				Description:  "Warm gooey dark chocolate cake with a molten center served with vanilla bean ice cream.",
				Price:        189.00,
				Category:     "Dessert",
				ImageURL:     "https://images.unsplash.com/photo-1606313564200-e75d5e30476c?w=600&auto=format&fit=crop&q=80",
				IsAvailable:  true,
			},
			{
				RestaurantID: restID,
				Name:         "Fresh Mango Smoothie",
				Description:  "Chilled tropical smoothie crafted with Alphonso mango pulp, Greek yogurt & honey.",
				Price:        149.00,
				Category:     "Drinks",
				ImageURL:     "https://images.unsplash.com/photo-1546173159-315724a31696?w=600&auto=format&fit=crop&q=80",
				IsAvailable:  true,
			},
		}

		for _, item := range items {
			db.Create(&item)
		}

		log.Println("Sample food items seeded successfully!")
	}
}
