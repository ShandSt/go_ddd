package product

import (
	"context"
	"time"

	"github.com/stasshander/ddd/internal/domain/product"
	"github.com/stasshander/ddd/internal/infrastructure/metrics"
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
	start := time.Now()

	p, err := product.NewProduct(name, description, price)
	if err != nil {
		metrics.ProductOperationsTotal.WithLabelValues("create", "validation_error").Inc()
		return nil, err
	}

	if err := s.repo.Create(ctx, p); err != nil {
		metrics.ProductOperationsTotal.WithLabelValues("create", "repository_error").Inc()
		return nil, err
	}

	duration := time.Since(start).Seconds()
	metrics.ProductOperationsTotal.WithLabelValues("create", "success").Inc()
	metrics.ProductOperationDuration.WithLabelValues("create").Observe(duration)

	return p, nil
}

func (s *Service) GetProduct(ctx context.Context, id string) (*product.Product, error) {
	start := time.Now()

	p, err := s.repo.GetByID(ctx, id)

	duration := time.Since(start).Seconds()
	status := "success"
	if err != nil {
		status = "error"
		if err == product.ErrProductNotFound {
			status = "not_found"
		}
	}

	metrics.ProductOperationsTotal.WithLabelValues("get", status).Inc()
	metrics.ProductOperationDuration.WithLabelValues("get").Observe(duration)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) UpdateProductPrice(ctx context.Context, id string, newPrice float64) error {
	start := time.Now()

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		metrics.ProductOperationsTotal.WithLabelValues("update_price", "not_found").Inc()
		return err
	}

	if err := p.UpdatePrice(newPrice); err != nil {
		metrics.ProductOperationsTotal.WithLabelValues("update_price", "validation_error").Inc()
		return err
	}

	if err := s.repo.Update(ctx, p); err != nil {
		metrics.ProductOperationsTotal.WithLabelValues("update_price", "repository_error").Inc()
		return err
	}

	duration := time.Since(start).Seconds()
	metrics.ProductOperationsTotal.WithLabelValues("update_price", "success").Inc()
	metrics.ProductOperationDuration.WithLabelValues("update_price").Observe(duration)

	return nil
}

func (s *Service) UpdateProductDescription(ctx context.Context, id string, newDescription string) error {
	start := time.Now()

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		metrics.ProductOperationsTotal.WithLabelValues("update_description", "not_found").Inc()
		return err
	}

	if err := p.UpdateDescription(newDescription); err != nil {
		metrics.ProductOperationsTotal.WithLabelValues("update_description", "validation_error").Inc()
		return err
	}

	if err := s.repo.Update(ctx, p); err != nil {
		metrics.ProductOperationsTotal.WithLabelValues("update_description", "repository_error").Inc()
		return err
	}

	duration := time.Since(start).Seconds()
	metrics.ProductOperationsTotal.WithLabelValues("update_description", "success").Inc()
	metrics.ProductOperationDuration.WithLabelValues("update_description").Observe(duration)

	return nil
}

func (s *Service) DeleteProduct(ctx context.Context, id string) error {
	start := time.Now()

	if err := s.repo.Delete(ctx, id); err != nil {
		metrics.ProductOperationsTotal.WithLabelValues("delete", "error").Inc()
		return err
	}

	duration := time.Since(start).Seconds()
	metrics.ProductOperationsTotal.WithLabelValues("delete", "success").Inc()
	metrics.ProductOperationDuration.WithLabelValues("delete").Observe(duration)

	return nil
}

func (s *Service) ListProducts(ctx context.Context) ([]*product.Product, error) {
	start := time.Now()

	products, err := s.repo.List(ctx)

	duration := time.Since(start).Seconds()
	status := "success"
	if err != nil {
		status = "error"
	}

	metrics.ProductOperationsTotal.WithLabelValues("list", status).Inc()
	metrics.ProductOperationDuration.WithLabelValues("list").Observe(duration)

	if err != nil {
		return nil, err
	}

	return products, nil
}
