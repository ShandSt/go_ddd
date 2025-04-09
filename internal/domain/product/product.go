package product

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Description Product entity
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewProduct creates a new product with the given parameters
func NewProduct(name, description string, price float64) *Product {
	now := time.Now()
	return &Product{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// UpdatePrice updates the product price
func (p *Product) UpdatePrice(price float64) {
	p.Price = price
	p.UpdatedAt = time.Now()
}

// UpdateDescription updates the product description
func (p *Product) UpdateDescription(description string) {
	p.Description = description
	p.UpdatedAt = time.Now()
}
