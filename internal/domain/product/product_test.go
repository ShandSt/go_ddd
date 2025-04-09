package product

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
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
			wantErr:     ErrInvalidPrice,
		},
		{
			name:        "empty name",
			productName: "",
			desc:        "Test Description",
			price:       10.0,
			wantErr:     ErrInvalidName,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := NewProduct(tc.productName, tc.desc, tc.price)
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
			assert.NotEmpty(t, p.ID)
			assert.NotEmpty(t, p.CreatedAt)
			assert.NotEmpty(t, p.UpdatedAt)
		})
	}
}

func TestUpdatePrice(t *testing.T) {
	testCases := []struct {
		name    string
		price   float64
		wantErr error
	}{
		{
			name:    "valid price",
			price:   20.0,
			wantErr: nil,
		},
		{
			name:    "zero price",
			price:   0.0,
			wantErr: ErrInvalidPrice,
		},
		{
			name:    "negative price",
			price:   -10.0,
			wantErr: ErrInvalidPrice,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := NewProduct("Test Product", "Test Description", 10.0)
			assert.NoError(t, err)

			err = p.UpdatePrice(tc.price)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Equal(t, 10.0, p.Price) // price should not change on error
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.price, p.Price)
		})
	}
}

func TestUpdateDescription(t *testing.T) {
	testCases := []struct {
		name        string
		description string
		wantErr     error
	}{
		{
			name:        "valid description",
			description: "New Description",
			wantErr:     nil,
		},
		{
			name:        "empty description",
			description: "",
			wantErr:     ErrInvalidDescription,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := NewProduct("Test Product", "Test Description", 10.0)
			assert.NoError(t, err)

			err = p.UpdateDescription(tc.description)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Equal(t, "Test Description", p.Description) // description should not change on error
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.description, p.Description)
		})
	}
}
