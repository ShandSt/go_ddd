package store

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, store *Store) error
	GetByID(ctx context.Context, id string) (*Store, error)
	Update(ctx context.Context, store *Store) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*Store, int, error)
	AddProduct(ctx context.Context, storeID string, productID string) error
	RemoveProduct(ctx context.Context, storeID string, productID string) error
}
