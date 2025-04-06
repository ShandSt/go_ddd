package product

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Description Product entity
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Price       float64            `bson:"price"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func NewProduct(name, description string, price float64) (*Product, error) {
	if name == "" {
		return nil, ErrInvalidName
	}

	if description == "" {
		return nil, ErrInvalidDescription
	}

	if price <= 0 {
		return nil, ErrInvalidPrice
	}

	now := time.Now()
	return &Product{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (p *Product) UpdatePrice(price float64) error {
	if price <= 0 {
		return ErrInvalidPrice
	}

	p.Price = price
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) UpdateDescription(description string) error {
	if description == "" {
		return ErrInvalidDescription
	}

	p.Description = description
	p.UpdatedAt = time.Now()
	return nil
}
