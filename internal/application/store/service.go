package store

import (
	"context"

	"github.com/stasshander/ddd/internal/domain/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Create(ctx context.Context, store *store.Store) error
	GetByID(ctx context.Context, id string) (*store.Store, error)
	UpdateName(ctx context.Context, id, name string) (*store.Store, error)
	UpdateAddress(ctx context.Context, id, address string) (*store.Store, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*store.Store, int, error)
	AddProduct(ctx context.Context, storeID string, productID primitive.ObjectID) (*store.Store, error)
	RemoveProduct(ctx context.Context, storeID string, productID primitive.ObjectID) (*store.Store, error)
}

type service struct {
	repo store.Repository
}

func NewService(repo store.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, store *store.Store) error {
	return s.repo.Create(ctx, store)
}

func (s *service) GetByID(ctx context.Context, id string) (*store.Store, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) UpdateName(ctx context.Context, id, name string) (*store.Store, error) {
	store, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	store.UpdateName(name)
	if err := s.repo.Update(ctx, store); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *service) UpdateAddress(ctx context.Context, id, address string) (*store.Store, error) {
	store, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	store.UpdateAddress(address)
	if err := s.repo.Update(ctx, store); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, page, limit int) ([]*store.Store, int, error) {
	return s.repo.List(ctx, page, limit)
}

func (s *service) AddProduct(ctx context.Context, storeID string, productID primitive.ObjectID) (*store.Store, error) {
	store, err := s.repo.GetByID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	store.AddProduct(productID)
	if err := s.repo.Update(ctx, store); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *service) RemoveProduct(ctx context.Context, storeID string, productID primitive.ObjectID) (*store.Store, error) {
	store, err := s.repo.GetByID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	store.RemoveProduct(productID)
	if err := s.repo.Update(ctx, store); err != nil {
		return nil, err
	}

	return store, nil
}
