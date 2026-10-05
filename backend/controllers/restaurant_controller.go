package controllers

import (
	"net/http"
	"strconv"

	"github.com/garvpatidar/food-delivery/backend/models"
	"github.com/garvpatidar/food-delivery/backend/usecases"
	"github.com/gin-gonic/gin"
)

type RestaurantController struct {
	restaurantUsecase usecases.RestaurantUsecase
}

func NewRestaurantController(u usecases.RestaurantUsecase) *RestaurantController {
	return &RestaurantController{
		restaurantUsecase: u,
	}
}

func (c *RestaurantController) Create(ctx *gin.Context) {
	var restaurant models.Restaurant

	// 1. Parse the incoming JSON
	if err := ctx.ShouldBindJSON(&restaurant); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 2. Pass to the Usecase (Business Logic)
	if err := c.restaurantUsecase.CreateRestaurant(ctx.Request.Context(), &restaurant); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Return success
	ctx.JSON(http.StatusCreated, gin.H{"message": "Restaurant created successfully", "restaurant": restaurant})
}

func (c *RestaurantController) Get(ctx *gin.Context) {
	// 1. Get ID from URL parameter (e.g. /restaurants/5)
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 2. Pass to Usecase
	restaurant, err := c.restaurantUsecase.GetRestaurant(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	// 3. Return response
	ctx.JSON(http.StatusOK, restaurant)
}

func (c *RestaurantController) GetAll(ctx *gin.Context) {
	restaurants, err := c.restaurantUsecase.GetAllRestaurants(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch restaurants"})
		return
	}

	ctx.JSON(http.StatusOK, restaurants)
}
