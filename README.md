# DDD Product Management API

A simple Domain-Driven Design (DDD) API for product management built with Go.

## Features

- Product management (CRUD operations)
- MongoDB integration
- Swagger API documentation
- Docker support
- Clean architecture with DDD principles

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- MongoDB

### Running with Docker

1. Build the Docker image:
```bash
docker build -t ddd-app .
```

2. Run MongoDB:
```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

3. Run the application:
```bash
docker run -d -p 8091:8080 --name ddd-container --link mongodb:mongodb -e MONGO_URI=mongodb://mongodb:27017 -e MONGO_DATABASE=products ddd-app
```

### API Documentation

Swagger UI is available at: http://localhost:8091/swagger/index.html

### API Endpoints

- `POST /api/products` - Create a new product
- `GET /api/products/:id` - Get a product by ID
- `PUT /api/products/:id/price` - Update product price
- `PUT /api/products/:id/description` - Update product description
- `DELETE /api/products/:id` - Delete a product
- `GET /api/products` - List all products

## Project Structure

```
.
├── cmd/
│   └── main.go
├── internal/
│   ├── application/
│   │   └── product/
│   ├── domain/
│   │   └── product/
│   ├── infrastructure/
│   │   ├── config/
│   ├── infrastructure/
│   │   └── mongodb/
│   └── interfaces/
│       └── http/
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── Dockerfile
├── go.mod
└── go.sum
```

## License

This project is licensed under the MIT License. 