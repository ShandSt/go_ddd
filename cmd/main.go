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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-contrib/cors"
	"github.com/stasshander/ddd/docs"
	"github.com/stasshander/ddd/internal/application/product"
	"github.com/stasshander/ddd/internal/infrastructure/config"
	"github.com/stasshander/ddd/internal/infrastructure/mongodb"
	httphandler "github.com/stasshander/ddd/internal/interfaces/http"
)

// @title           DDD API
// @version         1.0
// @description     A simple DDD API for product management
// @host            localhost:8080
// @BasePath        /api
func main() {
	cfg := config.LoadConfig()

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	db := client.Database(cfg.MongoDBName)
	repo := mongodb.NewProductRepository(db)
	service := product.NewService(repo)
	productHandler := httphandler.NewProductHandler(service)

	router := gin.Default()

	// Initialize Swagger
	docs.SwaggerInfo.Title = "DDD API"
	docs.SwaggerInfo.Description = "A simple DDD API for product management"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		products := api.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("/:id", productHandler.GetProduct)
			products.PUT("/:id/price", productHandler.UpdateProductPrice)
			products.PUT("/:id/description", productHandler.UpdateProductDescription)
			products.DELETE("/:id", productHandler.DeleteProduct)
			products.GET("", productHandler.ListProducts)
		}
	}

	srv := &http.Server{
		Addr:         cfg.BindAddress(),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
