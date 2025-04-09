package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appstore "github.com/stasshander/ddd/internal/application/store"
	domainstore "github.com/stasshander/ddd/internal/domain/store"
	"github.com/stasshander/ddd/internal/interfaces/http/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoreHandler struct {
	service *appstore.Service
}

func NewStoreHandler(service *appstore.Service) *StoreHandler {
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

	store, err := h.service.CreateStore(c.Request.Context(), req.Name, req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSimpleResponse(store))
}

func (h *StoreHandler) GetStore(c *gin.Context) {
	id := c.Param("id")
	store, err := h.service.GetStore(c.Request.Context(), id)
	if err != nil {
		if err == domainstore.ErrStoreNotFound {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSimpleResponse(store))
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

	if err := h.service.UpdateStoreName(c.Request.Context(), id, req.Name); err != nil {
		if err == domainstore.ErrStoreNotFound {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	store, err := h.service.GetStore(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSimpleResponse(store))
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

	if err := h.service.UpdateStoreAddress(c.Request.Context(), id, req.Address); err != nil {
		if err == domainstore.ErrStoreNotFound {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	store, err := h.service.GetStore(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSimpleResponse(store))
}

func (h *StoreHandler) DeleteStore(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteStore(c.Request.Context(), id); err != nil {
		if err == domainstore.ErrStoreNotFound {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSimpleResponse[any](nil))
}

func (h *StoreHandler) ListStores(c *gin.Context) {
	page := 1
	limit := 10

	stores, total, err := h.service.ListStores(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	pagination := &response.Pagination{
		Page:     page,
		PageSize: limit,
	}

	c.JSON(http.StatusOK, response.NewPaginatedResponse(stores, pagination, total))
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

	if err := h.service.AddProductToStore(c.Request.Context(), storeID, productID); err != nil {
		if err == domainstore.ErrStoreNotFound {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	store, err := h.service.GetStore(c.Request.Context(), storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSimpleResponse(store))
}

func (h *StoreHandler) RemoveProductFromStore(c *gin.Context) {
	storeID := c.Param("id")
	productIDStr := c.Param("productId")

	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(http.StatusBadRequest, "Invalid product ID"))
		return
	}

	if err := h.service.RemoveProductFromStore(c.Request.Context(), storeID, productID); err != nil {
		if err == domainstore.ErrStoreNotFound {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	store, err := h.service.GetStore(c.Request.Context(), storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSimpleResponse(store))
}
