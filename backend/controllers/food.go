
package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/garvpatidar/food-delivery/backend/database"
	"github.com/garvpatidar/food-delivery/backend/models"
	"github.com/gin-gonic/gin"
	redisClient "github.com/garvpatidar/food-delivery/backend/redis"
)

type CreateFoodInput struct {
	RestaurantID uint    `json:"restaurant_id"`
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description"`
	Price        float64 `json:"price" binding:"required,gt=0"`
	Category     string  `json:"category" binding:"required"`
	ImageURL     string  `json:"image_url"`
	IsAvailable  *bool   `json:"is_available"`
}

// GetFoodItems lists all available food items (or all items if all=true).
// Redis is checked first to avoid querying the database every time.
func GetFoodItems(c *gin.Context) {
	category := c.Query("category")
	showAll := c.Query("all") == "true"

	// Create a different cache key for each category.
	cacheKey := "food_items"
	if category != "" {
		cacheKey = "food_items:" + category
	}
	if showAll {
		cacheKey += ":all"
	}

	// Check Redis for cached food items if available.
	if redisClient.RDB != nil {
		cachedData, err := redisClient.RDB.Get(
			redisClient.Ctx,
			cacheKey,
		).Result()

		if err == nil {
			var items []models.FoodItem
			if err := json.Unmarshal([]byte(cachedData), &items); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"data":   items,
					"source": "redis",
				})
				return
			}
		}
	}

	// Database fetch
	var items []models.FoodItem
	query := database.DB

	if !showAll {
		query = query.Where("is_available = ?", true)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch food items",
		})
		return
	}

	// Store in Redis if available
	if redisClient.RDB != nil {
		data, err := json.Marshal(items)
		if err == nil {
			redisClient.RDB.Set(
				redisClient.Ctx,
				cacheKey,
				data,
				5*time.Minute,
			)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   items,
		"source": "database",
	})
}

// GetFoodItemByID retrieves a single food item by ID.
func GetFoodItemByID(c *gin.Context) {
	id := c.Param("id")
	var item models.FoodItem

	if err := database.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Food item not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

// CreateFoodItem creates a new food item.
func CreateFoodItem(c *gin.Context) {
	var input CreateFoodInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	restID := input.RestaurantID
	if restID == 0 {
		restID = 1
	}

	isAvailable := true
	if input.IsAvailable != nil {
		isAvailable = *input.IsAvailable
	}

	item := models.FoodItem{
		RestaurantID: restID,
		Name:         input.Name,
		Description:  input.Description,
		Price:        input.Price,
		Category:     input.Category,
		ImageURL:     input.ImageURL,
		IsAvailable:  isAvailable,
	}

	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create food item",
		})
		return
	}

	// Flush Redis food cache if connected
	if redisClient.RDB != nil {
		redisClient.RDB.FlushDB(redisClient.Ctx)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Food item created",
		"data":    item,
	})
}

// UpdateFoodItem updates an existing food item.
func UpdateFoodItem(c *gin.Context) {
	id := c.Param("id")
	var item models.FoodItem

	if err := database.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Food item not found",
		})
		return
	}

	var input CreateFoodInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	updates := map[string]interface{}{
		"name":        input.Name,
		"description": input.Description,
		"price":       input.Price,
		"category":    input.Category,
		"image_url":   input.ImageURL,
	}
	if input.IsAvailable != nil {
		updates["is_available"] = *input.IsAvailable
	}

	database.DB.Model(&item).Updates(updates)

	// Flush Redis food cache if connected
	if redisClient.RDB != nil {
		redisClient.RDB.FlushDB(redisClient.Ctx)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Food item updated",
		"data":    item,
	})
}

// DeleteFoodItem deletes a food item.
func DeleteFoodItem(c *gin.Context) {
	id := c.Param("id")
	var item models.FoodItem

	if err := database.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Food item not found",
		})
		return
	}

	database.DB.Delete(&item)

	// Delete food cache if connected
	if redisClient.RDB != nil {
		redisClient.RDB.FlushDB(redisClient.Ctx)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Food item deleted successfully",
	})
}

