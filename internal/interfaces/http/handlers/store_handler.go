package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stasshander/ddd/internal/application/store"
	"github.com/stasshander/ddd/internal/domain/store"
	"github.com/stasshander/ddd/internal/interfaces/http/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoreHandler struct {
	service store.Service
}

func NewStoreHandler(service store.Service) *StoreHandler {
	return &StoreHandler{
		service: service,
	}
}

type CreateStoreRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

func (h *StoreHandler) CreateStore(c *gin.Context) {
	var req CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(http.StatusBadRequest, "Invalid request body"))
		return
	}

	store := store.NewStore(req.Name, req.Address)
	if err := h.service.Create(c.Request.Context(), store); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSimpleResponse(http.StatusCreated, "Store created successfully", store))
}

func (h *StoreHandler) GetStore(c *gin.Context) {
	id := c.Param("id")
	store, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == store.ErrStoreNotFound {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSimpleResponse(http.StatusOK, "Store retrieved successfully", store))
}

type UpdateStoreNameRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *StoreHandler) UpdateStoreName(c *gin.Context) {
	id := c.Param("id")
	var req UpdateStoreNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(http.StatusBadRequest, "Invalid request body"))
		return
	}

	store, err := h.service.UpdateName(c.Request.Context(), id, req.Name)
	if err != nil {
		if err == store.ErrStoreNotFound {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSimpleResponse(http.StatusOK, "Store name updated successfully", store))
}

type UpdateStoreAddressRequest struct {
	Address string `json:"address" binding:"required"`
}

func (h *StoreHandler) UpdateStoreAddress(c *gin.Context) {
	id := c.Param("id")
	var req UpdateStoreAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(http.StatusBadRequest, "Invalid request body"))
		return
	}

	store, err := h.service.UpdateAddress(c.Request.Context(), id, req.Address)
	if err != nil {
		if err == store.ErrStoreNotFound {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSimpleResponse(http.StatusOK, "Store address updated successfully", store))
}

func (h *StoreHandler) DeleteStore(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err == store.ErrStoreNotFound {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSimpleResponse(http.StatusOK, "Store deleted successfully", nil))
}

func (h *StoreHandler) ListStores(c *gin.Context) {
	page := 1
	limit := 10

	stores, total, err := h.service.List(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewPaginatedResponse(http.StatusOK, "Stores retrieved successfully", stores, total, page, limit))
}

type AddProductRequest struct {
	ProductID string `json:"product_id" binding:"required"`
}

func (h *StoreHandler) AddProductToStore(c *gin.Context) {
	storeID := c.Param("id")
	var req AddProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(http.StatusBadRequest, "Invalid request body"))
		return
	}

	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(http.StatusBadRequest, "Invalid product ID"))
		return
	}

	store, err := h.service.AddProduct(c.Request.Context(), storeID, productID)
	if err != nil {
		if err == store.ErrStoreNotFound {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSimpleResponse(http.StatusOK, "Product added to store successfully", store))
}

func (h *StoreHandler) RemoveProductFromStore(c *gin.Context) {
	storeID := c.Param("id")
	productIDStr := c.Param("productId")

	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(http.StatusBadRequest, "Invalid product ID"))
		return
	}

	store, err := h.service.RemoveProduct(c.Request.Context(), storeID, productID)
	if err != nil {
		if err == store.ErrStoreNotFound {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSimpleResponse(http.StatusOK, "Product removed from store successfully", store))
}
