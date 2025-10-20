package main

import (
    "log"
    "os"

    "github.com/DavutcanJ/customer-segmentation-api/internal/config"
    "github.com/DavutcanJ/customer-segmentation-api/internal/handlers"
    "github.com/DavutcanJ/customer-segmentation-api/internal/middleware"
    "github.com/DavutcanJ/customer-segmentation-api/internal/services"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"

    _ "github.com/DavutcanJ/customer-segmentation-api/docs"
)

// @title Customer Segmentation API
// @version 1.0
// @description MiniBatchKMeans ve Ollama LLM ile müşteri segmentasyonu yapan Go API
// @host localhost:8008
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your Bearer token in the format: Bearer <token>
func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Connect to MongoDB
    config.ConnectDatabase()

    // Initialize Redis (optional, app works without it)
    if err := services.InitRedis(); err != nil {
        log.Printf("Warning: Redis connection failed: %v", err)
    }

    // Initialize database indexes
    userService := services.NewUserService()
    if err := userService.InitializeIndexes(); err != nil {
        log.Printf("Warning: Failed to create indexes: %v", err)
    }

    // Initialize Gin router
    r := gin.Default()

    // Swagger documentation
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Initialize handlers
    authHandler := handlers.NewAuthHandler()
    userHandler := handlers.NewUserHandler()
    segmentHandler := handlers.NewSegmentHandler()

    // Public routes
    api := r.Group("/api")
    {
        api.POST("/register", authHandler.Register)
        api.POST("/login", authHandler.Login)
    }

    // Protected routes
    protected := api.Group("/")
    protected.Use(middleware.JWTAuth())
    {
        protected.GET("/profile", userHandler.GetProfile)
        protected.POST("/segment", segmentHandler.SegmentCustomer)
    }

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok", "database": "connected"})
    })

    // Get port from environment or use default
    port := os.Getenv("PORT")
    if port == "" {
        port = "8008"
    }

    log.Printf("Server starting on port %s", port)
    log.Fatal(r.Run(":" + port))
}