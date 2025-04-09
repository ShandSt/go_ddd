package product

import (
	"context"
	"testing"

	"github.com/stasshander/ddd/internal/domain/product"
	"github.com/stretchr/testify/assert"
)

type MockRepository struct {
	products map[string]*product.Product
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		products: make(map[string]*product.Product),
	}
}

func (m *MockRepository) Create(ctx context.Context, p *product.Product) error {
	m.products[p.ID.Hex()] = p
	return nil
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*product.Product, error) {
	if p, ok := m.products[id]; ok {
		return p, nil
	}
	return nil, product.ErrProductNotFound
}

func (m *MockRepository) Update(ctx context.Context, p *product.Product) error {
	if _, ok := m.products[p.ID.Hex()]; !ok {
		return product.ErrProductNotFound
	}
	m.products[p.ID.Hex()] = p
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.products[id]; !ok {
		return product.ErrProductNotFound
	}
	delete(m.products, id)
	return nil
}

func (m *MockRepository) List(ctx context.Context) ([]*product.Product, error) {
	products := make([]*product.Product, 0, len(m.products))
	for _, p := range m.products {
		products = append(products, p)
	}
	return products, nil
}

func TestCreateProduct(t *testing.T) {
	testCases := []struct {
		name        string
		productName string
		desc        string
		price       float64
		wantErr     error
	}{
		{
			name:        "valid product",
			productName: "Test Product",
			desc:        "Test Description",
			price:       10.0,
			wantErr:     nil,
		},
		{
			name:        "invalid price",
			productName: "Test Product",
			desc:        "Test Description",
			price:       -10.0,
			wantErr:     product.ErrInvalidPrice,
		},
		{
			name:        "empty name",
			productName: "",
			desc:        "Test Description",
			price:       10.0,
			wantErr:     product.ErrInvalidName,
		},
		{
			name:        "empty description",
			productName: "Test Product",
			desc:        "",
			price:       10.0,
			wantErr:     product.ErrInvalidDescription,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockRepository()
			service := NewService(repo)

			p, err := service.CreateProduct(context.Background(), tc.productName, tc.desc, tc.price)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, p)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, p)
			assert.Equal(t, tc.productName, p.Name)
			assert.Equal(t, tc.desc, p.Description)
			assert.Equal(t, tc.price, p.Price)
		})
	}
}

func TestGetProduct(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(*Service) string
		wantErr error
	}{
		{
			name: "get existing product",
			setup: func(s *Service) string {
				p, _ := s.CreateProduct(context.Background(), "Test Product", "Test Description", 10.0)
				return p.ID.Hex()
			},
			wantErr: nil,
		},
		{
			name: "get non-existent product",
			setup: func(s *Service) string {
				return "nonexistentid"
			},
			wantErr: product.ErrProductNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockRepository()
			service := NewService(repo)

			id := tc.setup(service)
			p, err := service.GetProduct(context.Background(), id)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, p)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, p)
			assert.Equal(t, id, p.ID.Hex())
		})
	}
}

func TestUpdateProductPrice(t *testing.T) {
	testCases := []struct {
		name    string
		price   float64
		setup   func(*Service) string
		wantErr error
	}{
		{
			name:  "update price of existing product",
			price: 20.0,
			setup: func(s *Service) string {
				p, _ := s.CreateProduct(context.Background(), "Test Product", "Test Description", 10.0)
				return p.ID.Hex()
			},
			wantErr: nil,
		},
		{
			name:  "update price of non-existent product",
			price: 20.0,
			setup: func(s *Service) string {
				return "nonexistentid"
			},
			wantErr: product.ErrProductNotFound,
		},
		{
			name:  "update with invalid price",
			price: -20.0,
			setup: func(s *Service) string {
				p, _ := s.CreateProduct(context.Background(), "Test Product", "Test Description", 10.0)
				return p.ID.Hex()
			},
			wantErr: product.ErrInvalidPrice,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockRepository()
			service := NewService(repo)

			id := tc.setup(service)
			err := service.UpdateProductPrice(context.Background(), id, tc.price)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			assert.NoError(t, err)
			p, err := service.GetProduct(context.Background(), id)
			assert.NoError(t, err)
			assert.Equal(t, tc.price, p.Price)
		})
	}
}

func TestUpdateProductDescription(t *testing.T) {
	testCases := []struct {
		name        string
		description string
		setup       func(*Service) string
		wantErr     error
	}{
		{
			name:        "update description of existing product",
			description: "New Description",
			setup: func(s *Service) string {
				p, _ := s.CreateProduct(context.Background(), "Test Product", "Test Description", 10.0)
				return p.ID.Hex()
			},
			wantErr: nil,
		},
		{
			name:        "update description of non-existent product",
			description: "New Description",
			setup: func(s *Service) string {
				return "nonexistentid"
			},
			wantErr: product.ErrProductNotFound,
		},
		{
			name:        "update with empty description",
			description: "",
			setup: func(s *Service) string {
				p, _ := s.CreateProduct(context.Background(), "Test Product", "Test Description", 10.0)
				return p.ID.Hex()
			},
			wantErr: product.ErrInvalidDescription,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockRepository()
			service := NewService(repo)

			id := tc.setup(service)
			err := service.UpdateProductDescription(context.Background(), id, tc.description)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			assert.NoError(t, err)
			p, err := service.GetProduct(context.Background(), id)
			assert.NoError(t, err)
			assert.Equal(t, tc.description, p.Description)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(*Service) string
		wantErr error
	}{
		{
			name: "delete existing product",
			setup: func(s *Service) string {
				p, _ := s.CreateProduct(context.Background(), "Test Product", "Test Description", 10.0)
				return p.ID.Hex()
			},
			wantErr: nil,
		},
		{
			name: "delete non-existent product",
			setup: func(s *Service) string {
				return "nonexistentid"
			},
			wantErr: product.ErrProductNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockRepository()
			service := NewService(repo)

			id := tc.setup(service)
			err := service.DeleteProduct(context.Background(), id)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			assert.NoError(t, err)
			_, err = service.GetProduct(context.Background(), id)
			assert.ErrorIs(t, err, product.ErrProductNotFound)
		})
	}
}

func TestListProducts(t *testing.T) {
	testCases := []struct {
		name      string
		setup     func(*Service) int
		wantCount int
		wantErr   error
	}{
		{
			name: "list multiple products",
			setup: func(s *Service) int {
				_, err := s.CreateProduct(context.Background(), "Test Product 1", "Test Description 1", 10.0)
				if err != nil {
					return 0
				}
				_, err = s.CreateProduct(context.Background(), "Test Product 2", "Test Description 2", 20.0)
				if err != nil {
					return 1
				}
				return 2
			},
			wantCount: 2,
			wantErr:   nil,
		},
		{
			name: "list empty products",
			setup: func(s *Service) int {
				return 0
			},
			wantCount: 0,
			wantErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockRepository()
			service := NewService(repo)

			expectedCount := tc.setup(service)
			products, err := service.ListProducts(context.Background())

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, products)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, products, expectedCount)
		})
	}
}
