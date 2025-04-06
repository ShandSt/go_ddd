package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/stasshander/ddd/internal/domain/product"
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
	result, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	p.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*product.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var p product.Product
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
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": p.ID}, p)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *ProductRepository) List(ctx context.Context) ([]*product.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*product.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}
