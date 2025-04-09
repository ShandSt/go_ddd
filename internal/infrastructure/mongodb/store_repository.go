package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stasshander/ddd/internal/domain/store"
)

type StoreRepository struct {
	client       *mongo.Client
	databaseName string
	collection   *mongo.Collection
}

func NewStoreRepository(client *mongo.Client, databaseName string) *StoreRepository {
	collection := client.Database(databaseName).Collection("stores")
	return &StoreRepository{
		client:       client,
		databaseName: databaseName,
		collection:   collection,
	}
}

func (r *StoreRepository) Create(ctx context.Context, s *store.Store) error {
	result, err := r.collection.InsertOne(ctx, s)
	if err != nil {
		return err
	}

	s.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *StoreRepository) GetByID(ctx context.Context, id string) (*store.Store, error) {
	var s store.Store
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&s)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, store.ErrStoreNotFound
		}
		return nil, err
	}

	return &s, nil
}

func (r *StoreRepository) Update(ctx context.Context, s *store.Store) error {
	objectID, err := primitive.ObjectIDFromHex(s.ID.Hex())
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":       s.Name,
			"address":    s.Address,
			"products":   s.Products,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return store.ErrStoreNotFound
	}

	return nil
}

func (r *StoreRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return store.ErrStoreNotFound
	}

	return nil
}

func (r *StoreRepository) List(ctx context.Context, page, limit int) ([]*store.Store, int, error) {
	var stores []*store.Store

	// Calculate total count
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
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
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	if err = cursor.All(ctx, &stores); err != nil {
		return nil, 0, err
	}

	return stores, int(total), nil
}

func (r *StoreRepository) AddProduct(ctx context.Context, storeID string, productID string) error {
	storeObjectID, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return err
	}

	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": storeObjectID},
		bson.M{"$addToSet": bson.M{"product_ids": productObjectID}},
	)
	return err
}

func (r *StoreRepository) RemoveProduct(ctx context.Context, storeID string, productID string) error {
	storeObjectID, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return err
	}

	productObjectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": storeObjectID},
		bson.M{"$pull": bson.M{"product_ids": productObjectID}},
	)
	return err
}
