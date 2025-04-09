package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stasshander/ddd/internal/domain/product"
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
	result, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}

	p.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*product.Product, error) {
	var p product.Product
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *product.Product) error {
	objectID, err := primitive.ObjectIDFromHex(p.ID.Hex())
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

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return product.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) List(ctx context.Context) ([]*product.Product, error) {
	var products []*product.Product

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}
