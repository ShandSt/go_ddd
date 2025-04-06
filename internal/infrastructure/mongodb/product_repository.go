package mongodb

import (
	"context"
	// "time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/stasshander/ddd/internal/domain/product"
	// "github.com/stasshander/ddd/internal/infrastructure/metrics"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
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

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// metrics.MongoDBOperationsTotal.WithLabelValues("get_by_id", "error").Inc()
		return nil, err
	}

	var p product.Product
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

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": p.ID}, p)

	// duration := time.Since(start).Seconds()
	// status := "success"
	if err != nil {
		// status = "error"
		return err
	}

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

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})

	// duration := time.Since(start).Seconds()
	// status := "success"
	if err != nil {
		// status = "error"
		return err
	}

	// metrics.MongoDBOperationsTotal.WithLabelValues("delete", status).Inc()
	// metrics.MongoDBOperationDuration.WithLabelValues("delete").Observe(duration)

	return nil
}

func (r *ProductRepository) List(ctx context.Context) ([]*product.Product, error) {
	// start := time.Now()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		// metrics.MongoDBOperationsTotal.WithLabelValues("list", "error").Inc()
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*product.Product
	if err := cursor.All(ctx, &products); err != nil {
		// metrics.MongoDBOperationsTotal.WithLabelValues("list", "error").Inc()
		return nil, err
	}

	// duration := time.Since(start).Seconds()
	// metrics.MongoDBOperationsTotal.WithLabelValues("list", "success").Inc()
	// metrics.MongoDBOperationDuration.WithLabelValues("list").Observe(duration)

	return products, nil
}
