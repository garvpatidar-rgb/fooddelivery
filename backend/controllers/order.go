package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/garvpatidar/food-delivery/backend/database"
	"github.com/garvpatidar/food-delivery/backend/models"
	"github.com/garvpatidar/food-delivery/backend/rabbitmq"
	ws "github.com/garvpatidar/food-delivery/backend/websocket"
	"github.com/gin-gonic/gin"
)

// Hub is set by main.go so controllers can broadcast updates
var WSHub *ws.Hub

type OrderItemInput struct {
	FoodItemID uint `json:"food_item_id" binding:"required"`
	Quantity   int  `json:"quantity" binding:"required,gt=0"`
}

type CreateOrderInput struct {
	Items           []OrderItemInput `json:"items" binding:"required,gt=0"`
	DeliveryAddress string           `json:"delivery_address" binding:"required"`
	RestaurantID    uint             `json:"restaurant_id"` // Optional, defaults to 1
}

type UpdateOrderStatusInput struct {
	Status string `json:"status" binding:"required"`
}

// CreateOrder handles customer order placement
func CreateOrder(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDVal.(uint)

	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restID := input.RestaurantID
	if restID == 0 {
		restID = 1
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	// Calculate total amount and construct order items securely from DB
	for _, itemInput := range input.Items {
		var foodItem models.FoodItem
		if err := database.DB.First(&foodItem, itemInput.FoodItemID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food item ID"})
			return
		}

		itemTotal := foodItem.Price * float64(itemInput.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, models.OrderItem{
			FoodItemID: foodItem.ID,
			Quantity:   itemInput.Quantity,
			Price:      foodItem.Price,
		})
	}

	order := models.Order{
		UserID:          userID,
		RestaurantID:    restID,
		Items:           orderItems,
		TotalAmount:     totalAmount,
		Status:          "pending",
		DeliveryAddress: input.DeliveryAddress,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Preload items for full response
	database.DB.Preload("Items.FoodItem").First(&order, order.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order placed successfully",
		"data":    order,
	})
}

// GetUserOrders fetches all orders for the logged in user
func GetUserOrders(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)

	var orders []models.Order
	if err := database.DB.Preload("Items.FoodItem").Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// GetAllOrders fetches all orders in the system for Admin / Restaurant Dashboard
func GetAllOrders(c *gin.Context) {
	var orders []models.Order
	if err := database.DB.Preload("Items.FoodItem").Preload("User").Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// GetOrderByID fetches details of a specific order
func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)

	var order models.Order
	if err := database.DB.Preload("Items.FoodItem").Preload("User").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Security check: ensure user owns the order or is admin
	roleVal, _ := c.Get("userRole")
	role := roleVal.(string)

	if order.UserID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// UpdateOrderStatus updates order status, publishes RabbitMQ event, and broadcasts WebSocket update
func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var input UpdateOrderStatusInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = input.Status
	database.DB.Save(&order)

	// 1. Publish to RabbitMQ (if available)
	event := rabbitmq.OrderStatusEvent{
		OrderID: order.ID,
		UserID:  order.UserID,
		Status:  order.Status,
		Message: "Order status updated to: " + order.Status,
	}
	if err := rabbitmq.PublishOrderStatusEvent(context.Background(), event); err != nil {
		log.Printf("Note: RabbitMQ event skip: %v", err)
	}

	// 2. Broadcast real-time update via WebSocket directly to the user
	if WSHub != nil {
		BroadcastOrderUpdate(WSHub, order.UserID, order.ID, order.Status)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated", "data": order})
}
