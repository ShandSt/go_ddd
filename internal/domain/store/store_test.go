package store

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewStore(t *testing.T) {
	tests := []struct {
		name      string
		storeName string
		address   string
		wantErr   bool
		errType   error
	}{
		{
			name:      "Valid store",
			storeName: "Test Store",
			address:   "123 Test St",
			wantErr:   false,
		},
		{
			name:      "Empty name",
			storeName: "",
			address:   "123 Test St",
			wantErr:   true,
			errType:   ErrInvalidStoreName,
		},
		{
			name:      "Empty address",
			storeName: "Test Store",
			address:   "",
			wantErr:   true,
			errType:   ErrInvalidStoreAddress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, err := NewStore(tt.storeName, tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err != tt.errType {
					t.Errorf("NewStore() error = %v, want error %v", err, tt.errType)
				}
				return
			}
			if store == nil {
				t.Error("NewStore() returned nil store")
				return
			}
			if store.Name != tt.storeName {
				t.Errorf("NewStore() store.Name = %v, want %v", store.Name, tt.storeName)
			}
			if store.Address != tt.address {
				t.Errorf("NewStore() store.Address = %v, want %v", store.Address, tt.address)
			}
			if store.ProductIDs == nil {
				t.Error("NewStore() store.ProductIDs is nil")
			}
			if len(store.ProductIDs) != 0 {
				t.Errorf("NewStore() store.ProductIDs length = %v, want 0", len(store.ProductIDs))
			}
			if store.CreatedAt.IsZero() {
				t.Error("NewStore() store.CreatedAt is zero")
			}
			if store.UpdatedAt.IsZero() {
				t.Error("NewStore() store.UpdatedAt is zero")
			}
		})
	}
}

func TestStore_AddProduct(t *testing.T) {
	store, _ := NewStore("Test Store", "123 Test St")
	productID := primitive.NewObjectID()

	tests := []struct {
		name      string
		productID primitive.ObjectID
		wantErr   bool
		errType   error
	}{
		{
			name:      "Add new product",
			productID: productID,
			wantErr:   false,
		},
		{
			name:      "Add existing product",
			productID: productID,
			wantErr:   true,
			errType:   ErrProductAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.AddProduct(tt.productID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.AddProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err != tt.errType {
					t.Errorf("Store.AddProduct() error = %v, want error %v", err, tt.errType)
				}
				return
			}
			if len(store.ProductIDs) != 1 {
				t.Errorf("Store.AddProduct() store.ProductIDs length = %v, want 1", len(store.ProductIDs))
			}
			if store.ProductIDs[0] != tt.productID {
				t.Errorf("Store.AddProduct() store.ProductIDs[0] = %v, want %v", store.ProductIDs[0], tt.productID)
			}
		})
	}
}

func TestStore_RemoveProduct(t *testing.T) {
	store, _ := NewStore("Test Store", "123 Test St")
	productID := primitive.NewObjectID()
	store.AddProduct(productID)

	tests := []struct {
		name      string
		productID primitive.ObjectID
		wantErr   bool
		errType   error
	}{
		{
			name:      "Remove existing product",
			productID: productID,
			wantErr:   false,
		},
		{
			name:      "Remove non-existing product",
			productID: primitive.NewObjectID(),
			wantErr:   true,
			errType:   ErrProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.RemoveProduct(tt.productID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.RemoveProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err != tt.errType {
					t.Errorf("Store.RemoveProduct() error = %v, want error %v", err, tt.errType)
				}
				return
			}
			if len(store.ProductIDs) != 0 {
				t.Errorf("Store.RemoveProduct() store.ProductIDs length = %v, want 0", len(store.ProductIDs))
			}
		})
	}
}

func TestStore_HasProduct(t *testing.T) {
	store, _ := NewStore("Test Store", "123 Test St")
	productID := primitive.NewObjectID()
	store.AddProduct(productID)

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

func TestStore_UpdateName(t *testing.T) {
	store, _ := NewStore("Test Store", "123 Test St")

	tests := []struct {
		name    string
		newName string
		wantErr bool
		errType error
	}{
		{
			name:    "Valid name",
			newName: "New Store Name",
			wantErr: false,
		},
		{
			name:    "Empty name",
			newName: "",
			wantErr: true,
			errType: ErrInvalidStoreName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := store.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure UpdatedAt will be different

			err := store.UpdateName(tt.newName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.UpdateName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err != tt.errType {
					t.Errorf("Store.UpdateName() error = %v, want error %v", err, tt.errType)
				}
				return
			}
			if store.Name != tt.newName {
				t.Errorf("Store.UpdateName() store.Name = %v, want %v", store.Name, tt.newName)
			}
			if !store.UpdatedAt.After(oldUpdatedAt) {
				t.Errorf("Store.UpdateName() store.UpdatedAt = %v, want after %v", store.UpdatedAt, oldUpdatedAt)
			}
		})
	}
}

func TestStore_UpdateAddress(t *testing.T) {
	store, _ := NewStore("Test Store", "123 Test St")

	tests := []struct {
		name       string
		newAddress string
		wantErr    bool
		errType    error
	}{
		{
			name:       "Valid address",
			newAddress: "456 New St",
			wantErr:    false,
		},
		{
			name:       "Empty address",
			newAddress: "",
			wantErr:    true,
			errType:    ErrInvalidStoreAddress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := store.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure UpdatedAt will be different

			err := store.UpdateAddress(tt.newAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.UpdateAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err != tt.errType {
					t.Errorf("Store.UpdateAddress() error = %v, want error %v", err, tt.errType)
				}
				return
			}
			if store.Address != tt.newAddress {
				t.Errorf("Store.UpdateAddress() store.Address = %v, want %v", store.Address, tt.newAddress)
			}
			if !store.UpdatedAt.After(oldUpdatedAt) {
				t.Errorf("Store.UpdateAddress() store.UpdatedAt = %v, want after %v", store.UpdatedAt, oldUpdatedAt)
			}
		})
	}
}
