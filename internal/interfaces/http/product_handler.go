package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appProduct "github.com/stasshander/ddd/internal/application/product"
	"github.com/stasshander/ddd/internal/domain/product"
)

type ProductHandler struct {
	service *appProduct.Service
}

func NewProductHandler(service *appProduct.Service) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) RegisterRoutes(router *gin.Engine) {
	products := router.Group("/api/products")
	{
		products.POST("/", h.CreateProduct)
		products.GET("/:id", h.GetProduct)
		products.PUT("/:id/price", h.UpdateProductPrice)
		products.PUT("/:id/description", h.UpdateProductDescription)
		products.DELETE("/:id", h.DeleteProduct)
		products.GET("/", h.ListProducts)
	}
}

// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateProductRequest true "Product details"
// @Success 201 {object} ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	p, err := h.service.CreateProduct(c.Request.Context(), req.Name, req.Description, req.Price)
	if err != nil {
		if err == product.ErrInvalidPrice {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toProductResponse(p))
}

// @Summary Get a product by ID
// @Description Get a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} ProductResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	p, err := h.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		if err == product.ErrProductNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, toProductResponse(p))
}

// @Summary Update product price
// @Description Update the price of a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param price body UpdatePriceRequest true "New price"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id}/price [put]
func (h *ProductHandler) UpdateProductPrice(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.UpdateProductPrice(c.Request.Context(), id, req.Price); err != nil {
		if err == product.ErrProductNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Product not found"})
			return
		}
		if err == product.ErrInvalidPrice {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Update product description
// @Description Update the description of a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param description body UpdateDescriptionRequest true "New description"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id}/description [put]
func (h *ProductHandler) UpdateProductDescription(c *gin.Context) {
	id := c.Param("id")
	var req UpdateDescriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.UpdateProductDescription(c.Request.Context(), id, req.Description); err != nil {
		if err == product.ErrProductNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		if err == product.ErrProductNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary List all products
// @Description Get a list of all products
// @Tags products
// @Produce json
// @Success 200 {array} ProductResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.service.ListProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := make([]ProductResponse, len(products))
	for i, p := range products {
		response[i] = toProductResponse(p)
	}

	c.JSON(http.StatusOK, response)
}

func toProductResponse(p *product.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
