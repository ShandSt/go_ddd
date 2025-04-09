package store

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidStoreName     = errors.New("store name cannot be empty")
	ErrInvalidStoreAddress  = errors.New("store address cannot be empty")
	ErrStoreNotFound        = errors.New("store not found")
	ErrProductAlreadyExists = errors.New("product already exists in store")
	ErrProductNotFound      = errors.New("product not found in store")
)

type Store struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name      string               `bson:"name" json:"name"`
	Address   string               `bson:"address" json:"address"`
	Products  []primitive.ObjectID `bson:"products" json:"products"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

func NewStore(name, address string) (*Store, error) {
	if name == "" {
		return nil, ErrInvalidStoreName
	}
	if address == "" {
		return nil, ErrInvalidStoreAddress
	}

	now := time.Now()
	return &Store{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Address:   address,
		Products:  make([]primitive.ObjectID, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *Store) AddProduct(productID primitive.ObjectID) error {
	for _, id := range s.Products {
		if id == productID {
			return ErrProductAlreadyExists
		}
	}
	s.Products = append(s.Products, productID)
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Store) RemoveProduct(productID primitive.ObjectID) error {
	for i, id := range s.Products {
		if id == productID {
			s.Products = append(s.Products[:i], s.Products[i+1:]...)
			s.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrProductNotFound
}

func (s *Store) HasProduct(productID primitive.ObjectID) bool {
	for _, id := range s.Products {
		if id == productID {
			return true
		}
	}
	return false
}

func (s *Store) UpdateName(name string) error {
	if name == "" {
		return ErrInvalidStoreName
	}
	s.Name = name
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Store) UpdateAddress(address string) error {
	if address == "" {
		return ErrInvalidStoreAddress
	}
	s.Address = address
	s.UpdatedAt = time.Now()
	return nil
}
