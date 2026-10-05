package routes

import (
	"net/http"

	"github.com/garvpatidar/food-delivery/backend/controllers"
	"github.com/garvpatidar/food-delivery/backend/database"
	"github.com/garvpatidar/food-delivery/backend/middleware"
	redisClient "github.com/garvpatidar/food-delivery/backend/redis"
	"github.com/garvpatidar/food-delivery/backend/repositories"
	"github.com/garvpatidar/food-delivery/backend/usecases"
	ws "github.com/garvpatidar/food-delivery/backend/websocket"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, hub *ws.Hub) {
	// Enable CORS
	router.Use(middleware.CORSMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Food Delivery API is running",
		})
	})

	// ─── Dependency Injection ─────────────────────────────────────────
	// Restaurant
	restaurantRepo := repositories.NewRestaurantRepository(database.DB)
	restaurantUsecase := usecases.NewRestaurantUsecase(restaurantRepo)
	restaurantController := controllers.NewRestaurantController(restaurantUsecase)

	// Cart (Redis)
	cartRepo := repositories.NewCartRepository(redisClient.RDB)
	cartUsecase := usecases.NewCartUsecase(cartRepo)
	cartController := controllers.NewCartController(cartUsecase)

	// Payment
	paymentRepo := repositories.NewPaymentRepository(database.DB)
	paymentUsecase := usecases.NewPaymentUsecase(paymentRepo)
	paymentController := controllers.NewPaymentController(paymentUsecase)

	// WebSocket
	wsController := controllers.NewWebSocketController(hub)
	// ─────────────────────────────────────────────────────────────────

	api := router.Group("/api/v1")
	{
		// Public Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// Public Food Routes
		food := api.Group("/food")
		{
			food.GET("", controllers.GetFoodItems)
			food.GET("/:id", controllers.GetFoodItemByID)
		}

		// Public Restaurant Routes (read-only)
		restaurants := api.Group("/restaurants")
		{
			restaurants.GET("", restaurantController.GetAll)
			restaurants.GET("/:id", restaurantController.Get)
		}

		// Protected User Routes
		protected := api.Group("/user")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", func(c *gin.Context) {
				userID, _ := c.Get("userID")
				email, _ := c.Get("userEmail")
				role, _ := c.Get("userRole")

				c.JSON(http.StatusOK, gin.H{
					"user_id": userID,
					"email":   email,
					"role":    role,
				})
			})
		}

		// Protected Order Routes (Customer)
		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.POST("", controllers.CreateOrder)
			orders.GET("", controllers.GetUserOrders)
			orders.GET("/:id", controllers.GetOrderByID)
			orders.PUT("/:id/status", controllers.UpdateOrderStatus)
		}

		// Cart Routes (Protected)
		cart := api.Group("/cart")
		cart.Use(middleware.AuthMiddleware())
		{
			cart.GET("", cartController.GetCart)
			cart.POST("/items", cartController.AddItem)
			cart.DELETE("/items/:foodItemId", cartController.RemoveItem)
			cart.DELETE("", cartController.ClearCart)
		}

		// Payment Routes (Protected)
		payments := api.Group("/payments")
		payments.Use(middleware.AuthMiddleware())
		{
			payments.POST("", paymentController.InitiatePayment)
			payments.PUT("/:id/confirm", paymentController.ConfirmPayment)
			payments.GET("/order/:orderID", paymentController.GetPaymentByOrder)
		}

		// Admin Food Routes
		adminFood := api.Group("/admin/food")
		adminFood.Use(middleware.AuthMiddleware())
		{
			adminFood.POST("", controllers.CreateFoodItem)
			adminFood.PUT("/:id", controllers.UpdateFoodItem)
			adminFood.DELETE("/:id", controllers.DeleteFoodItem)
		}

		// Admin Restaurant Management
		adminRestaurants := api.Group("/admin/restaurants")
		adminRestaurants.Use(middleware.AuthMiddleware())
		{
			adminRestaurants.POST("", restaurantController.Create)
		}

		// Admin Order Management Routes
		adminOrders := api.Group("/admin/orders")
		adminOrders.Use(middleware.AuthMiddleware())
		{
			adminOrders.GET("", controllers.GetAllOrders)
			adminOrders.PUT("/:id/status", controllers.UpdateOrderStatus)
		}
	}

	// WebSocket Route (Protected) - Real-time order tracking
	router.GET("/ws", middleware.AuthMiddleware(), wsController.ServeWS)
}
