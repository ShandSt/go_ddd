# DDD Product Store API

A Go-based REST API for managing products and stores, built using Domain-Driven Design principles.

## Features

- Product management (CRUD operations)
- Store management with product associations
- MongoDB for data persistence
- Prometheus metrics for monitoring
- Clean architecture with DDD principles
- Docker support
- API token authentication

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- MongoDB

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/stasshander/ddd.git
cd ddd
```

2. Start the MongoDB container:
```bash
docker-compose up -d mongodb
```

3. Build and run the application:
```bash
docker-compose up --build app
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | HTTP server port | 8091 |
| SERVER_HOST | HTTP server host | localhost |
| READ_TIMEOUT | Server read timeout | 10s |
| WRITE_TIMEOUT | Server write timeout | 10s |
| IDLE_TIMEOUT | Server idle timeout | 60s |
| READ_HEADER_TIMEOUT | Server read header timeout | 2s |
| MONGO_URI | MongoDB connection string | mongodb://localhost:27017 |
| MONGO_DATABASE | MongoDB database name | products |
| API_TOKEN | API authentication token | "" |

## API Endpoints

### Products

- `GET /api/products` - List all products
- `POST /api/products` - Create a new product
- `GET /api/products/:id` - Get product by ID
- `PUT /api/products/:id/price` - Update product price
- `PUT /api/products/:id/description` - Update product description
- `DELETE /api/products/:id` - Delete product

### Stores

- `GET /api/stores` - List all stores
- `POST /api/stores` - Create a new store
- `GET /api/stores/:id` - Get store by ID
- `PUT /api/stores/:id/name` - Update store name
- `PUT /api/stores/:id/address` - Update store address
- `DELETE /api/stores/:id` - Delete store
- `POST /api/stores/:id/products` - Add product to store
- `DELETE /api/stores/:id/products/:productId` - Remove product from store

### Metrics

- `GET /metrics` - Prometheus metrics endpoint

## Project Structure

```
.
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── application/      # Application services
│   ├── domain/          # Domain models and interfaces
│   ├── infrastructure/  # External implementations
│   └── interfaces/      # HTTP handlers and middleware
├── docker-compose.yml
└── README.md
```

## Testing

Run the tests:
```bash
go test ./... -v
```

## Metrics

The application exposes Prometheus metrics at the `/metrics` endpoint, including:
- HTTP request counts and durations
- Go runtime metrics
- MongoDB operation metrics

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 