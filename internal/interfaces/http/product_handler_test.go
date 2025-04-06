package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	appProduct "github.com/stasshander/ddd/internal/application/product"
	"github.com/stasshander/ddd/internal/domain/product"
	"github.com/stretchr/testify/assert"
)

type MockProductRepository struct {
	products map[string]*product.Product
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		products: make(map[string]*product.Product),
	}
}

func (m *MockProductRepository) Create(ctx context.Context, p *product.Product) error {
	m.products[p.ID.Hex()] = p
	return nil
}

func (m *MockProductRepository) GetByID(ctx context.Context, id string) (*product.Product, error) {
	if p, ok := m.products[id]; ok {
		return p, nil
	}
	return nil, product.ErrProductNotFound
}

func (m *MockProductRepository) Update(ctx context.Context, p *product.Product) error {
	if _, ok := m.products[p.ID.Hex()]; !ok {
		return product.ErrProductNotFound
	}
	m.products[p.ID.Hex()] = p
	return nil
}

func (m *MockProductRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.products[id]; !ok {
		return product.ErrProductNotFound
	}
	delete(m.products, id)
	return nil
}

func (m *MockProductRepository) List(ctx context.Context) ([]*product.Product, error) {
	products := make([]*product.Product, 0, len(m.products))
	for _, p := range m.products {
		products = append(products, p)
	}
	return products, nil
}

func setupTest() (*gin.Engine, *ProductHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := NewMockProductRepository()
	service := appProduct.NewService(repo)
	handler := NewProductHandler(service)
	handler.RegisterRoutes(router)

	return router, handler
}

func createTestProduct(t *testing.T, handler *ProductHandler) *product.Product {
	p, err := product.NewProduct("Test Product", "Test Description", 10.0)
	assert.NoError(t, err)

	createdProduct, err := handler.service.CreateProduct(context.Background(), p.Name, p.Description, p.Price)
	assert.NoError(t, err)
	return createdProduct
}

func TestCreateProduct(t *testing.T) {
	testCases := []struct {
		name       string
		request    CreateProductRequest
		wantStatus int
		wantError  bool
	}{
		{
			name: "create valid product",
			request: CreateProductRequest{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       10.0,
			},
			wantStatus: http.StatusCreated,
			wantError:  false,
		},
		{
			name: "create product with invalid price",
			request: CreateProductRequest{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       -10.0,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router, _ := setupTest()

			body, err := json.Marshal(tc.request)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/products/", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			if !tc.wantError {
				var response ProductResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.request.Name, response.Name)
				assert.Equal(t, tc.request.Description, response.Description)
				assert.Equal(t, tc.request.Price, response.Price)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	testCases := []struct {
		name       string
		setup      func(*ProductHandler) string
		wantStatus int
		wantError  bool
	}{
		{
			name: "get existing product",
			setup: func(h *ProductHandler) string {
				p := createTestProduct(t, h)
				return p.ID.Hex()
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name: "get non-existent product",
			setup: func(h *ProductHandler) string {
				return "nonexistentid"
			},
			wantStatus: http.StatusNotFound,
			wantError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router, handler := setupTest()
			productID := tc.setup(handler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/products/"+productID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			if !tc.wantError {
				var response ProductResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.ID)
				assert.NotEmpty(t, response.Name)
				assert.NotEmpty(t, response.Description)
				assert.Greater(t, response.Price, float64(0))
			}
		})
	}
}

func TestUpdateProductPrice(t *testing.T) {
	testCases := []struct {
		name       string
		setup      func(*ProductHandler) string
		price      float64
		wantStatus int
		wantError  bool
	}{
		{
			name: "update price of existing product",
			setup: func(h *ProductHandler) string {
				p := createTestProduct(t, h)
				return p.ID.Hex()
			},
			price:      20.0,
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name: "update price of non-existent product",
			setup: func(h *ProductHandler) string {
				return "nonexistentid"
			},
			price:      20.0,
			wantStatus: http.StatusNotFound,
			wantError:  true,
		},
		{
			name: "update price with invalid value",
			setup: func(h *ProductHandler) string {
				p := createTestProduct(t, h)
				return p.ID.Hex()
			},
			price:      -20.0,
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router, handler := setupTest()
			productID := tc.setup(handler)

			reqBody := UpdatePriceRequest{Price: tc.price}
			body, err := json.Marshal(reqBody)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/api/products/"+productID+"/price", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			if !tc.wantError {
				updatedProduct, err := handler.service.GetProduct(context.Background(), productID)
				assert.NoError(t, err)
				assert.Equal(t, tc.price, updatedProduct.Price)
			}
		})
	}
}

func TestUpdateProductDescription(t *testing.T) {
	testCases := []struct {
		name        string
		setup       func(*ProductHandler) string
		description string
		wantStatus  int
		wantError   bool
	}{
		{
			name: "update description of existing product",
			setup: func(h *ProductHandler) string {
				p := createTestProduct(t, h)
				return p.ID.Hex()
			},
			description: "New Description",
			wantStatus:  http.StatusOK,
			wantError:   false,
		},
		{
			name: "update description of non-existent product",
			setup: func(h *ProductHandler) string {
				return "nonexistentid"
			},
			description: "New Description",
			wantStatus:  http.StatusNotFound,
			wantError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router, handler := setupTest()
			productID := tc.setup(handler)

			reqBody := UpdateDescriptionRequest{Description: tc.description}
			body, err := json.Marshal(reqBody)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/api/products/"+productID+"/description", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			if !tc.wantError {
				updatedProduct, err := handler.service.GetProduct(context.Background(), productID)
				assert.NoError(t, err)
				assert.Equal(t, tc.description, updatedProduct.Description)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	testCases := []struct {
		name       string
		setup      func(*ProductHandler) string
		wantStatus int
		wantError  bool
	}{
		{
			name: "delete existing product",
			setup: func(h *ProductHandler) string {
				p := createTestProduct(t, h)
				return p.ID.Hex()
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name: "delete non-existent product",
			setup: func(h *ProductHandler) string {
				return "nonexistentid"
			},
			wantStatus: http.StatusNotFound,
			wantError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router, handler := setupTest()
			productID := tc.setup(handler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/products/"+productID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			if !tc.wantError {
				_, err := handler.service.GetProduct(context.Background(), productID)
				assert.ErrorIs(t, err, product.ErrProductNotFound)
			}
		})
	}
}

func TestListProducts(t *testing.T) {
	testCases := []struct {
		name       string
		setup      func(*ProductHandler) int
		wantStatus int
		wantCount  int
		wantError  bool
	}{
		{
			name: "list multiple products",
			setup: func(h *ProductHandler) int {
				createTestProduct(t, h)
				createTestProduct(t, h)
				return 2
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
			wantError:  false,
		},
		{
			name: "list empty products",
			setup: func(h *ProductHandler) int {
				return 0
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
			wantError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router, handler := setupTest()
			expectedCount := tc.setup(handler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/products/", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			if !tc.wantError {
				var response []ProductResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, expectedCount, len(response))
			}
		})
	}
}
