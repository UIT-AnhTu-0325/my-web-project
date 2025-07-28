package main

import (
	"log"
	"os"

	"hotel-backend/internal/database"
	"hotel-backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database connection
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", getEnv("FRONTEND_URL", "http://localhost:3000"))
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Hotel E-commerce API is running",
		})
	})

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	roomHandler := handlers.NewRoomHandler(db)
	productHandler := handlers.NewProductHandler(db)
	cartHandler := handlers.NewCartHandler(db)
	orderHandler := handlers.NewOrderHandler(db)
	adminHandler := handlers.NewAdminHandler(db)

	// API routes
	api := r.Group("/api")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/send-otp", authHandler.SendOTP)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/profile", authHandler.GetProfile)
		}

		// Rooms routes
		rooms := api.Group("/rooms")
		{
			rooms.GET("", roomHandler.GetRooms)
			rooms.GET("/:id", roomHandler.GetRoomByID)
			rooms.POST("/check-availability", roomHandler.CheckRoomAvailability)
		}

		// Products routes
		products := api.Group("/products")
		{
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProductByID)
			products.GET("/categories", productHandler.GetProductCategories)
		}

		// Cart routes
		cart := api.Group("/cart")
		{
			cart.GET("", cartHandler.GetCartItems)
			cart.POST("/add", cartHandler.AddToCart)
			cart.DELETE("/:id", cartHandler.RemoveFromCart)
			cart.DELETE("/clear", cartHandler.ClearCart)
		}

		// Orders routes
		orders := api.Group("/orders")
		{
			orders.GET("", orderHandler.GetOrders)
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrderByID)
		}

		// Admin routes
		admin := api.Group("/admin")
		{
			admin.GET("/dashboard", adminHandler.GetDashboardStats)
			admin.GET("/orders", adminHandler.GetAllOrders)
			admin.PUT("/orders/:id", adminHandler.UpdateOrderStatus)
			admin.POST("/rooms", adminHandler.AddRoom)
			admin.POST("/products", adminHandler.AddProduct)
		}
	}

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
