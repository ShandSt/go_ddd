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
	_ "github.com/stasshander/ddd/docs"
	"github.com/stasshander/ddd/internal/application/product"
	"github.com/stasshander/ddd/internal/application/store"
	"github.com/stasshander/ddd/internal/infrastructure/config"
	"github.com/stasshander/ddd/internal/infrastructure/metrics"
	"github.com/stasshander/ddd/internal/infrastructure/mongodb"
	"github.com/stasshander/ddd/internal/interfaces/http/handlers"
	"github.com/stasshander/ddd/internal/interfaces/http/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	productRepo := mongodb.NewProductRepository(client, cfg.MongoDB.Database)
	storeRepo := mongodb.NewStoreRepository(client, cfg.MongoDB.Database)

	productService := product.NewService(productRepo)
	storeService := store.NewService(storeRepo)

	router := gin.Default()

	router.Use(middleware.MetricsMiddleware())
	router.Use(middleware.AuthMiddleware(cfg.API.Token))

	productHandler := handlers.NewProductHandler(productService)
	storeHandler := handlers.NewStoreHandler(storeService)

	api := router.Group("/api")
	{
		products := api.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.ListProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.PUT("/:id/price", productHandler.UpdateProductPrice)
			products.PUT("/:id/description", productHandler.UpdateProductDescription)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}

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

	router.GET("/metrics", gin.WrapH(metrics.Handler()))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	srv := &http.Server{
		Addr:              cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:           router,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
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
