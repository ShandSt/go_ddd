package http

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required" example:"Product Name"`
	Description string  `json:"description" binding:"required" example:"Product Description"`
	Price       float64 `json:"price" binding:"required" example:"99.99"`
}

type UpdatePriceRequest struct {
	Price float64 `json:"price" binding:"required" example:"149.99"`
}

type UpdateDescriptionRequest struct {
	Description string `json:"description" binding:"required" example:"Updated Product Description"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}

type ProductResponse struct {
	ID          string  `json:"id" example:"507f1f77bcf86cd799439011"`
	Name        string  `json:"name" example:"Product Name"`
	Description string  `json:"description" example:"Product Description"`
	Price       float64 `json:"price" example:"99.99"`
	CreatedAt   string  `json:"created_at" example:"2024-04-06T11:22:31Z"`
	UpdatedAt   string  `json:"updated_at" example:"2024-04-06T11:22:31Z"`
}
