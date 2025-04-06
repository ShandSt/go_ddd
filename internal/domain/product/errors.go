package product

import "errors"

var (
	ErrInvalidPrice       = errors.New("invalid price")
	ErrProductNotFound    = errors.New("product not found")
	ErrInvalidName        = errors.New("invalid name")
	ErrInvalidDescription = errors.New("invalid description")
)
