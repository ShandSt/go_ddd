package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewStore(t *testing.T) {
	name := "Test Store"
	address := "123 Test St"

	store, err := NewStore(name, address)
	assert.NoError(t, err)
	assert.NotNil(t, store)
	assert.Equal(t, name, store.Name)
	assert.Equal(t, address, store.Address)
	assert.NotEmpty(t, store.ID)
	assert.NotZero(t, store.CreatedAt)
	assert.NotZero(t, store.UpdatedAt)
}

func TestUpdateName(t *testing.T) {
	store, err := NewStore("Old Name", "123 Test St")
	assert.NoError(t, err)

	newName := "New Name"
	err = store.UpdateName(newName)
	assert.NoError(t, err)
	assert.Equal(t, newName, store.Name)
	assert.True(t, store.UpdatedAt.After(store.CreatedAt))
}

func TestUpdateAddress(t *testing.T) {
	store, err := NewStore("Test Store", "Old Address")
	assert.NoError(t, err)

	newAddress := "New Address"
	err = store.UpdateAddress(newAddress)
	assert.NoError(t, err)
	assert.Equal(t, newAddress, store.Address)
	assert.True(t, store.UpdatedAt.After(store.CreatedAt))
}

func TestAddProduct(t *testing.T) {
	store, err := NewStore("Test Store", "Test Address")
	assert.NoError(t, err)

	productID := primitive.NewObjectID()
	err = store.AddProduct(productID)
	assert.NoError(t, err)
	assert.Contains(t, store.Products, productID)
	assert.True(t, store.UpdatedAt.After(store.CreatedAt))

	// Test adding the same product again
	err = store.AddProduct(productID)
	assert.Error(t, err)
	assert.Equal(t, ErrProductAlreadyExists, err)
}

func TestRemoveProduct(t *testing.T) {
	store, err := NewStore("Test Store", "123 Test St")
	assert.NoError(t, err)

	productID := primitive.NewObjectID()
	err = store.AddProduct(productID)
	assert.NoError(t, err)

	err = store.RemoveProduct(productID)
	assert.NoError(t, err)
	assert.NotContains(t, store.Products, productID)
	assert.True(t, store.UpdatedAt.After(store.CreatedAt))

	// Test removing non-existent product
	err = store.RemoveProduct(productID)
	assert.Error(t, err)
	assert.Equal(t, ErrProductNotFound, err)
}

func TestStore_HasProduct(t *testing.T) {
	store, err := NewStore("Test Store", "123 Test St")
	assert.NoError(t, err)

	productID := primitive.NewObjectID()
	err = store.AddProduct(productID)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		productID primitive.ObjectID
		want      bool
	}{
		{
			name:      "Existing product",
			productID: productID,
			want:      true,
		},
		{
			name:      "Non-existing product",
			productID: primitive.NewObjectID(),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := store.HasProduct(tt.productID)
			if got != tt.want {
				t.Errorf("Store.HasProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}
