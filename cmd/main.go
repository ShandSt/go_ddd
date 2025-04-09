package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stasshander/ddd/internal/application/product"
	"github.com/stasshander/ddd/internal/application/store"
	"github.com/stasshander/ddd/internal/infrastructure/config"
	"github.com/stasshander/ddd/internal/infrastructure/metrics"
	"github.com/stasshander/ddd/internal/infrastructure/mongodb"
	_ "github.com/stasshander/ddd/internal/interfaces/http/docs"
	"github.com/stasshander/ddd/internal/interfaces/http/handlers"
	"github.com/stasshander/ddd/internal/interfaces/http/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title DDD Product Store API
// @version 1.0
// @description A Go-based REST API for managing products and stores, built using Domain-Driven Design principles.
// @host localhost:8091
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize MongoDB client
	client, err := mongodb.NewClient(cfg.MongoDB.URI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Initialize repositories
	productRepo := mongodb.NewProductRepository(client, cfg.MongoDB.Database)
	storeRepo := mongodb.NewStoreRepository(client, cfg.MongoDB.Database)

	// Initialize services
	productService := product.NewService(productRepo)
	storeService := store.NewService(storeRepo)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.MetricsMiddleware())
	router.Use(middleware.AuthMiddleware(cfg.API.Token))

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)
	storeHandler := handlers.NewStoreHandler(storeService)

	// API routes
	api := router.Group("/api")
	{
		// Product routes
		products := api.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.ListProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.PUT("/:id/price", productHandler.UpdateProductPrice)
			products.PUT("/:id/description", productHandler.UpdateProductDescription)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}

		// Store routes
		stores := api.Group("/stores")
		{
			stores.POST("", storeHandler.CreateStore)
			stores.GET("", storeHandler.ListStores)
			stores.GET("/:id", storeHandler.GetStore)
			stores.PUT("/:id/name", storeHandler.UpdateStoreName)
			stores.PUT("/:id/address", storeHandler.UpdateStoreAddress)
			stores.DELETE("/:id", storeHandler.DeleteStore)
			stores.POST("/:id/products", storeHandler.AddProductToStore)
			stores.DELETE("/:id/products/:productId", storeHandler.RemoveProductFromStore)
		}
	}

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(metrics.Handler()))

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create HTTP server
	srv := &http.Server{
		Addr:              cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:           router,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
