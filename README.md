# DDD Project

A simple Domain-Driven Design (DDD) project in Go.

## Features

- Clean Architecture
- Domain-Driven Design
- MongoDB Integration
- RESTful API
- Swagger Documentation
- Docker Support

## Prerequisites

- Go 1.21 or higher
- MongoDB
- Docker (optional)

## Installation

1. Clone the repository:
```bash
git clone git@github.com:ShandSt/ddd.git
cd ddd
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables:
```bash
cp .env.example .env
```

4. Run the application:
```bash
make run
```

## Docker

Build and run with Docker:
```bash
make docker-build
make docker-run
```

## API Documentation

Swagger documentation is available at:
```
http://localhost:8080/swagger/index.html
```

## Environment Variables

- SERVER_PORT: Server port (default: 8080)
- SERVER_HOST: Server host (default: localhost)
- READ_TIMEOUT: Read timeout (default: 10s)
- WRITE_TIMEOUT: Write timeout (default: 10s)
- IDLE_TIMEOUT: Idle timeout (default: 120s)
- READ_HEADER_TIMEOUT: Read header timeout (default: 5s)
- MONGO_URI: MongoDB connection URI (default: mongodb://localhost:27017)
- MONGO_DATABASE: MongoDB database name (default: ddd)
- API_TOKEN: API authentication token (default: changethis)
- BIND_ADDRESS: Server bind address (default: localhost:8080)

## Testing

Run tests:
```bash
make test
```

## Linting

Run linter:
```bash
make lint
``` 