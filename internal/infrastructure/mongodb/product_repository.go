package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stasshander/ddd/internal/domain/product"
	// "github.com/stasshander/ddd/internal/infrastructure/metrics"
)

type ProductRepository struct {
	client       *mongo.Client
	databaseName string
	collection   *mongo.Collection
}

func NewProductRepository(client *mongo.Client, databaseName string) *ProductRepository {
	collection := client.Database(databaseName).Collection("products")
	return &ProductRepository{
		client:       client,
		databaseName: databaseName,
		collection:   collection,
	}
}

func (r *ProductRepository) Create(ctx context.Context, p *product.Product) error {
	// start := time.Now()

	result, err := r.collection.InsertOne(ctx, p)

	// duration := time.Since(start).Seconds()
	// status := "success"
	if err != nil {
		// status = "error"
		return err
	}

	p.ID = result.InsertedID.(primitive.ObjectID)

	// metrics.MongoDBOperationsTotal.WithLabelValues("create", status).Inc()
	// metrics.MongoDBOperationDuration.WithLabelValues("create").Observe(duration)

	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*product.Product, error) {
	// start := time.Now()

	var p product.Product
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// metrics.MongoDBOperationsTotal.WithLabelValues("get_by_id", "error").Inc()
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&p)

	// duration := time.Since(start).Seconds()
	// status := "success"
	if err != nil {
		// status = "error"
		if err == mongo.ErrNoDocuments {
			// metrics.MongoDBOperationsTotal.WithLabelValues("get_by_id", "not_found").Inc()
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}

	// metrics.MongoDBOperationsTotal.WithLabelValues("get_by_id", status).Inc()
	// metrics.MongoDBOperationDuration.WithLabelValues("get_by_id").Observe(duration)

	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *product.Product) error {
	// start := time.Now()

	objectID, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":        p.Name,
			"description": p.Description,
			"price":       p.Price,
			"updated_at":  time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return product.ErrProductNotFound
	}

	// duration := time.Since(start).Seconds()
	// status := "success"
	// metrics.MongoDBOperationsTotal.WithLabelValues("update", status).Inc()
	// metrics.MongoDBOperationDuration.WithLabelValues("update").Observe(duration)

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	// start := time.Now()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// metrics.MongoDBOperationsTotal.WithLabelValues("delete", "error").Inc()
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		// status = "error"
		return err
	}

	if result.DeletedCount == 0 {
		// metrics.MongoDBOperationsTotal.WithLabelValues("delete", "not_found").Inc()
		return product.ErrProductNotFound
	}

	// duration := time.Since(start).Seconds()
	// status := "success"
	// metrics.MongoDBOperationsTotal.WithLabelValues("delete", status).Inc()
	// metrics.MongoDBOperationDuration.WithLabelValues("delete").Observe(duration)

	return nil
}

func (r *ProductRepository) List(ctx context.Context, page, limit int) ([]*product.Product, int, error) {
	// start := time.Now()

	var products []*product.Product

	// Calculate total count
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		// metrics.MongoDBOperationsTotal.WithLabelValues("list", "error").Inc()
		return nil, 0, err
	}

	// Set up pagination options
	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Execute query
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		// metrics.MongoDBOperationsTotal.WithLabelValues("list", "error").Inc()
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	if err = cursor.All(ctx, &products); err != nil {
		// metrics.MongoDBOperationsTotal.WithLabelValues("list", "error").Inc()
		return nil, 0, err
	}

	// duration := time.Since(start).Seconds()
	// metrics.MongoDBOperationsTotal.WithLabelValues("list", "success").Inc()
	// metrics.MongoDBOperationDuration.WithLabelValues("list").Observe(duration)

	return products, int(total), nil
}
