package product

import "errors"

// Common product domain errors
var (
	// ErrInvalidPrice is returned when a product price is invalid (less than or equal to 0)
	ErrInvalidPrice = errors.New("invalid price")

	// ErrProductNotFound is returned when a requested product cannot be found
	ErrProductNotFound = errors.New("product not found")

	// ErrInvalidName is returned when a product name is invalid (empty)
	ErrInvalidName = errors.New("invalid name")

	// ErrInvalidDescription is returned when a product description is invalid (empty)
	ErrInvalidDescription = errors.New("invalid description")
)
