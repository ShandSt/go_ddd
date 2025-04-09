package store

import (
	"context"

	"github.com/stasshander/ddd/internal/domain/store"
	"github.com/stasshander/ddd/internal/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo *mongodb.StoreRepository
}

func NewService(repo *mongodb.StoreRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateStore(ctx context.Context, name, address string) (*store.Store, error) {
	store, err := store.NewStore(name, address)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, store); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *Service) GetStore(ctx context.Context, id string) (*store.Store, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateStoreName(ctx context.Context, id, name string) error {
	store, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := store.UpdateName(name); err != nil {
		return err
	}

	return s.repo.Update(ctx, store)
}

func (s *Service) UpdateStoreAddress(ctx context.Context, id, address string) error {
	store, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := store.UpdateAddress(address); err != nil {
		return err
	}

	return s.repo.Update(ctx, store)
}

func (s *Service) DeleteStore(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) ListStores(ctx context.Context, page, limit int) ([]*store.Store, int, error) {
	return s.repo.List(ctx, page, limit)
}

func (s *Service) AddProductToStore(ctx context.Context, storeID string, productID primitive.ObjectID) error {
	store, err := s.repo.GetByID(ctx, storeID)
	if err != nil {
		return err
	}

	if err := store.AddProduct(productID); err != nil {
		return err
	}

	return s.repo.Update(ctx, store)
}

func (s *Service) RemoveProductFromStore(ctx context.Context, storeID string, productID primitive.ObjectID) error {
	store, err := s.repo.GetByID(ctx, storeID)
	if err != nil {
		return err
	}

	if err := store.RemoveProduct(productID); err != nil {
		return err
	}

	return s.repo.Update(ctx, store)
}
