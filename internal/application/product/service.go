package product

import (
	"context"

	"github.com/stasshander/ddd/internal/domain/product"
)

type Service struct {
	repo product.Repository
}

func NewService(repo product.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateProduct(ctx context.Context, name, description string, price float64) (*product.Product, error) {
	p, err := product.NewProduct(name, description, price)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) GetProduct(ctx context.Context, id string) (*product.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) UpdateProductPrice(ctx context.Context, id string, newPrice float64) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := p.UpdatePrice(newPrice); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateProductDescription(ctx context.Context, id string, newDescription string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := p.UpdateDescription(newDescription); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProduct(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *Service) ListProducts(ctx context.Context) ([]*product.Product, error) {
	products, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}
