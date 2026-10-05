package main

import (
	"log"

	"github.com/garvpatidar/food-delivery/backend/config"
	"github.com/garvpatidar/food-delivery/backend/controllers"
	"github.com/garvpatidar/food-delivery/backend/database"
	"github.com/garvpatidar/food-delivery/backend/rabbitmq"
	"github.com/garvpatidar/food-delivery/backend/routes"
	ws "github.com/garvpatidar/food-delivery/backend/websocket"
	"github.com/gin-gonic/gin"
	redisClient "github.com/garvpatidar/food-delivery/backend/redis"
)

func main() {
	// Load Configuration
	cfg := config.Load()

	// Connect to Database & Run Auto-Migrations
	database.ConnectDB(cfg)

	// Connect to Redis
	redisClient.ConnectRedis(
		cfg.RedisAddr,
		cfg.RedisUsername,
		cfg.RedisPassword,
	)

	// Connect to RabbitMQ (non-fatal if not available in dev)
	if cfg.RabbitMQURL != "" {
		if err := rabbitmq.ConnectRabbitMQ(cfg.RabbitMQURL); err != nil {
			log.Printf("Warning: RabbitMQ connection failed: %v (continuing without it)", err)
		} else {
			defer rabbitmq.Close()
		}
	}

	// Initialize WebSocket Hub and run it in background
	hub := ws.NewHub()
	go hub.Run()

	// Make hub accessible to order controller for real-time broadcasts
	controllers.WSHub = hub

	// Initialize Gin Router
	router := gin.Default()

	// Setup API Routes (pass hub for WebSocket wiring)
	routes.SetupRoutes(router, hub)

	// Start Server
	log.Println("Server starting on :8080")
	router.Run(":8080")
}