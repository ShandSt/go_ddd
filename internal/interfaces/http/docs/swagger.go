// Package docs provides Swagger documentation for the API
package docs

import "github.com/swaggo/swag"

// @title DDD Product Store API
// @version 1.0
// @description A Go-based REST API for managing products and stores, built using Domain-Driven Design principles.
// @host localhost:8091
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @tag.name products
// @tag.description Product management endpoints

// @tag.name stores
// @tag.description Store management endpoints

// Product represents a product in the system
// @Description Product model
type Product struct {
	ID          string  `json:"id" example:"5f9d7b3b9d3f2b1b3c9d3f2b"`
	Name        string  `json:"name" example:"Product Name"`
	Description string  `json:"description" example:"Product Description"`
	Price       float64 `json:"price" example:"99.99"`
	CreatedAt   string  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   string  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// Store represents a store in the system
// @Description Store model
type Store struct {
	ID        string   `json:"id" example:"5f9d7b3b9d3f2b1b3c9d3f2b"`
	Name      string   `json:"name" example:"Store Name"`
	Address   string   `json:"address" example:"123 Store Street"`
	Products  []string `json:"products" example:"['5f9d7b3b9d3f2b1b3c9d3f2b']"`
	CreatedAt string   `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt string   `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// SimpleResponse represents a simple response structure
// @Description Simple response model
type SimpleResponse struct {
	Success bool        `json:"success" example:"true"`
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data"`
}

// PaginatedResponse represents a paginated response structure
// @Description Paginated response model
type PaginatedResponse struct {
	Success bool        `json:"success" example:"true"`
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data"`
	Total   int         `json:"total" example:"100"`
	Page    int         `json:"page" example:"1"`
	Limit   int         `json:"limit" example:"10"`
}

// CreateProductRequest represents the request to create a product
// @Description Create product request model
type CreateProductRequest struct {
	Name        string  `json:"name" example:"New Product"`
	Description string  `json:"description" example:"Product Description"`
	Price       float64 `json:"price" example:"99.99"`
}

// UpdateProductPriceRequest represents the request to update a product's price
// @Description Update product price request model
type UpdateProductPriceRequest struct {
	Price float64 `json:"price" example:"149.99"`
}

// UpdateProductDescriptionRequest represents the request to update a product's description
// @Description Update product description request model
type UpdateProductDescriptionRequest struct {
	Description string `json:"description" example:"Updated product description"`
}

// CreateStoreRequest represents the request to create a store
// @Description Create store request model
type CreateStoreRequest struct {
	Name    string `json:"name" example:"New Store"`
	Address string `json:"address" example:"456 New Street"`
}

// UpdateStoreNameRequest represents the request to update a store's name
// @Description Update store name request model
type UpdateStoreNameRequest struct {
	Name string `json:"name" example:"Updated Store Name"`
}

// UpdateStoreAddressRequest represents the request to update a store's address
// @Description Update store address request model
type UpdateStoreAddressRequest struct {
	Address string `json:"address" example:"789 Updated Street"`
}

// AddProductToStoreRequest represents the request to add a product to a store
// @Description Add product to store request model
type AddProductToStoreRequest struct {
	ProductID string `json:"product_id" example:"5f9d7b3b9d3f2b1b3c9d3f2b"`
}

func init() {
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  docTemplate,
	})
}

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/products": {
            "get": {
                "security": [{"ApiKeyAuth": []}],
                "description": "List all products",
                "produces": ["application/json"],
                "tags": ["products"],
                "summary": "List all products",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/PaginatedResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [{"ApiKeyAuth": []}],
                "description": "Create a new product",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["products"],
                "summary": "Create a new product",
                "parameters": [{
                    "description": "Product data",
                    "name": "product",
                    "in": "body",
                    "required": true,
                    "schema": {
                        "$ref": "#/definitions/CreateProductRequest"
                    }
                }],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/SimpleResponse"
                        }
                    }
                }
            }
        },
        "/products/{id}": {
            "get": {
                "security": [{"ApiKeyAuth": []}],
                "description": "Get product by ID",
                "produces": ["application/json"],
                "tags": ["products"],
                "summary": "Get product by ID",
                "parameters": [{
                    "type": "string",
                    "description": "Product ID",
                    "name": "id",
                    "in": "path",
                    "required": true
                }],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SimpleResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [{"ApiKeyAuth": []}],
                "description": "Delete product",
                "produces": ["application/json"],
                "tags": ["products"],
                "summary": "Delete product",
                "parameters": [{
                    "type": "string",
                    "description": "Product ID",
                    "name": "id",
                    "in": "path",
                    "required": true
                }],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SimpleResponse"
                        }
                    }
                }
            }
        },
        "/products/{id}/price": {
            "put": {
                "security": [{"ApiKeyAuth": []}],
                "description": "Update product price",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["products"],
                "summary": "Update product price",
                "parameters": [{
                    "type": "string",
                    "description": "Product ID",
                    "name": "id",
                    "in": "path",
                    "required": true
                }, {
                    "description": "New price",
                    "name": "price",
                    "in": "body",
                    "required": true,
                    "schema": {
                        "$ref": "#/definitions/UpdateProductPriceRequest"
                    }
                }],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SimpleResponse"
                        }
                    }
                }
            }
        },
        "/products/{id}/description": {
            "put": {
                "security": [{"ApiKeyAuth": []}],
                "description": "Update product description",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["products"],
                "summary": "Update product description",
                "parameters": [{
                    "type": "string",
                    "description": "Product ID",
                    "name": "id",
                    "in": "path",
                    "required": true
                }, {
                    "description": "New description",
                    "name": "description",
                    "in": "body",
                    "required": true,
                    "schema": {
                        "$ref": "#/definitions/UpdateProductDescriptionRequest"
                    }
                }],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SimpleResponse"
                        }
                    }
                }
            }
        },
        "/stores": {
            "get": {
                "security": [{"ApiKeyAuth": []}],
                "description": "List all stores",
                "produces": ["application/json"],
                "tags": ["stores"],
                "summary": "List all stores",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/PaginatedResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [{"ApiKeyAuth": []}],
                "description": "Create a new store",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["stores"],
                "summary": "Create a new store",
                "parameters": [{
                    "description": "Store data",
                    "name": "store",
                    "in": "body",
                    "required": true,
                    "schema": {
                        "$ref": "#/definitions/CreateStoreRequest"
                    }
                }],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/SimpleResponse"
                        }
                    }
                }
            }
        },
        "/stores/{id}": {
            "get": {
                "security": [{"ApiKeyAuth": []}],
                "description": "Get store by ID",
                "produces": ["application/json"],
                "tags": ["stores"],
                "summary": "Get store by ID",
                "parameters": [{
                    "type": "string",
                    "description": "Store ID",
                    "name": "id",
                    "in": "path",
                    "required": true
                }],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SimpleResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [{"ApiKeyAuth": []}],
                "description": "Delete store",
                "produces": ["application/json"],
                "tags": ["stores"],
                "summary": "Delete store",
                "parameters": [{
                    "type": "string",
                    "description": "Store ID",
                    "name": "id",
                    "in": "path",
                    "required": true
                }],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SimpleResponse"
                        }
                    }
                }
            }
        },
        "/stores/{id}/name": {
            "put": {
                "security": [{"ApiKeyAuth": []}],
                "description": "Update store name",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["stores"],
                "summary": "Update store name",
                "parameters": [{
                    "type": "string",
                    "description": "Store ID",
                    "name": "id",
                    "in": "path",
                    "required": true
                }, {
                    "description": "New name",
                    "name": "name",
                    "in": "body",
                    "required": true,
                    "schema": {
                        "$ref": "#/definitions/UpdateStoreNameRequest"
                    }
                }],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SimpleResponse"
                        }
                    }
                }
            }
        },
        "/stores/{id}/address": {
            "put": {
                "security": [{"ApiKeyAuth": []}],
                "description": "Update store address",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["stores"],
                "summary": "Update store address",
                "parameters": [{
                    "type": "string",
                    "description": "Store ID",
                    "name": "id",
                    "in": "path",
                    "required": true
                }, {
                    "description": "New address",
                    "name": "address",
                    "in": "body",
                    "required": true,
                    "schema": {
                        "$ref": "#/definitions/UpdateStoreAddressRequest"
                    }
                }],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SimpleResponse"
                        }
                    }
                }
            }
        },
        "/stores/{id}/products": {
            "post": {
                "security": [{"ApiKeyAuth": []}],
                "description": "Add product to store",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["stores"],
                "summary": "Add product to store",
                "parameters": [{
                    "type": "string",
                    "description": "Store ID",
                    "name": "id",
                    "in": "path",
                    "required": true
                }, {
                    "description": "Product ID",
                    "name": "product",
                    "in": "body",
                    "required": true,
                    "schema": {
                        "$ref": "#/definitions/AddProductToStoreRequest"
                    }
                }],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SimpleResponse"
                        }
                    }
                }
            }
        },
        "/stores/{id}/products/{productId}": {
            "delete": {
                "security": [{"ApiKeyAuth": []}],
                "description": "Remove product from store",
                "produces": ["application/json"],
                "tags": ["stores"],
                "summary": "Remove product from store",
                "parameters": [{
                    "type": "string",
                    "description": "Store ID",
                    "name": "id",
                    "in": "path",
                    "required": true
                }, {
                    "type": "string",
                    "description": "Product ID",
                    "name": "productId",
                    "in": "path",
                    "required": true
                }],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SimpleResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "Product": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "5f9d7b3b9d3f2b1b3c9d3f2b"
                },
                "name": {
                    "type": "string",
                    "example": "Product Name"
                },
                "description": {
                    "type": "string",
                    "example": "Product Description"
                },
                "price": {
                    "type": "number",
                    "format": "float",
                    "example": 99.99
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2023-01-01T00:00:00Z"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2023-01-01T00:00:00Z"
                }
            }
        },
        "Store": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "5f9d7b3b9d3f2b1b3c9d3f2b"
                },
                "name": {
                    "type": "string",
                    "example": "Store Name"
                },
                "address": {
                    "type": "string",
                    "example": "123 Store Street"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": ["5f9d7b3b9d3f2b1b3c9d3f2b"]
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2023-01-01T00:00:00Z"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2023-01-01T00:00:00Z"
                }
            }
        },
        "SimpleResponse": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "boolean",
                    "example": true
                },
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "message": {
                    "type": "string",
                    "example": "Operation successful"
                },
                "data": {
                    "type": "object"
                }
            }
        },
        "PaginatedResponse": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "boolean",
                    "example": true
                },
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "message": {
                    "type": "string",
                    "example": "Operation successful"
                },
                "data": {
                    "type": "object"
                },
                "total": {
                    "type": "integer",
                    "example": 100
                },
                "page": {
                    "type": "integer",
                    "example": 1
                },
                "limit": {
                    "type": "integer",
                    "example": 10
                }
            }
        },
        "CreateProductRequest": {
            "type": "object",
            "required": ["name", "description", "price"],
            "properties": {
                "name": {
                    "type": "string",
                    "example": "New Product"
                },
                "description": {
                    "type": "string",
                    "example": "Product Description"
                },
                "price": {
                    "type": "number",
                    "format": "float",
                    "example": 99.99
                }
            }
        },
        "UpdateProductPriceRequest": {
            "type": "object",
            "required": ["price"],
            "properties": {
                "price": {
                    "type": "number",
                    "format": "float",
                    "example": 149.99
                }
            }
        },
        "UpdateProductDescriptionRequest": {
            "type": "object",
            "required": ["description"],
            "properties": {
                "description": {
                    "type": "string",
                    "example": "Updated product description"
                }
            }
        },
        "CreateStoreRequest": {
            "type": "object",
            "required": ["name", "address"],
            "properties": {
                "name": {
                    "type": "string",
                    "example": "New Store"
                },
                "address": {
                    "type": "string",
                    "example": "456 New Street"
                }
            }
        },
        "UpdateStoreNameRequest": {
            "type": "object",
            "required": ["name"],
            "properties": {
                "name": {
                    "type": "string",
                    "example": "Updated Store Name"
                }
            }
        },
        "UpdateStoreAddressRequest": {
            "type": "object",
            "required": ["address"],
            "properties": {
                "address": {
                    "type": "string",
                    "example": "789 Updated Street"
                }
            }
        },
        "AddProductToStoreRequest": {
            "type": "object",
            "required": ["product_id"],
            "properties": {
                "product_id": {
                    "type": "string",
                    "example": "5f9d7b3b9d3f2b1b3c9d3f2b"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`
