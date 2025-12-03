# PetShop API

This repository hosts PetShop API — a modular backend service implemented in Go that provides REST APIs for pet shop management operations. The project follows a Ports & Adapters (Hexagonal) architecture and is designed for testability, maintainability, and scalability.

---

## Overview

PetShop API is a management system focused on helping pet shops handle customer information, addresses, phone contacts, and scheduling operations. The API provides endpoints for customer management, address management, and related business operations with robust validation and caching mechanisms.

### Technology stack
- Go 1.25
- chi router for HTTP
- GORM for PostgreSQL integration
- Redis for caching
- Kafka (franz-go) for event streaming
- Docker & Docker Compose for local environments
- zap for structured logging

### Key features
- Customer management with CPF/CNPJ validation
- Address management with validation
- Phone contact management (mobile and landline)
- Redis caching for improved performance
- Event-driven architecture with Kafka support
- Comprehensive test coverage with mocks
- Database migrations and seed data
- Brazilian document validation (CPF, CNPJ, area codes)

---

## Installation

### Prerequisites
- Go 1.25 or newer
- Git
- Docker & Docker Compose (for containerized development)
- Make (GNU Make)

Clone and prepare
```bash
git clone git@github.com:lechitz/petshop-api.git
cd petshop-api
```

Download modules
```bash
go mod download
```

---

## Configuration

### Environment variables

The application uses environment variables for configuration. Key variables include:

**Server configuration**
```bash
SERVER_CONTEXT=petshop-api     # API context path
PORT=5001                      # Server port
READ_TIMEOUT=10s               # HTTP read timeout
WRITE_TIMEOUT=10s              # HTTP write timeout
```

**Database configuration**
```bash
DB_USER=petshop-system         # PostgreSQL username
DB_PASSWORD=test1234           # PostgreSQL password
DB_NAME=petshop-system         # Database name
DB_HOST=localhost              # Database host
DB_PORT=5432                   # Database port
DB_TYPE=postgres               # Database type
```

**Redis configuration**
```bash
REDIS_ADDR=localhost:6379      # Redis address
REDIS_PASSWORD=                # Redis password (optional)
REDIS_DB=0                     # Redis database number
POOL_SIZE=100                  # Connection pool size
```

**Kafka configuration** (optional)
```bash
KAFKA_SCHEDULE_BOOTSTRAP_SERVER=localhost:29092
KAFKA_SCHEDULE_GROUPID=kafka_schedule
KAFKA_SCHEDULE_AUTO_OFFSET_RESET=earliest
KAFKA_SCHEDULE_TOPIC=schedule
```

### Start development environment

Start all services with Docker Compose
```bash
make docker-compose-up
```

This will start:
- PostgreSQL database on port 5432
- Redis cache on port 6379
- API Gateway on port 9999

For development mode (without rebuilding)
```bash
make docker-compose-dev-up
```

---

## Development & common commands

### Build the application
```bash
go build -o bin/petshop-api ./cmd/petshop-api
```

### Run the application locally
```bash
export DB_USER=petshop-system
export DB_PASSWORD=test1234
export DB_NAME=petshop-system
export DB_HOST=localhost
export DB_PORT=5432
export REDIS_ADDR=localhost:6379
./bin/petshop-api
```

### Run tests
```bash
go test ./...
```

### Run tests with coverage
```bash
make test-cover
```
This generates a coverage report excluding mock files and opens it in your browser.

### Docker commands

Build Docker image
```bash
make docker-build
```

Build and run Docker container
```bash
make docker-build-run
```

Run Docker container (image must exist)
```bash
make docker-run
```

Stop and clean up Docker environment
```bash
make docker-compose-down
```

---

## API summary

Base URL: `http://localhost:5001/petshop-api`

### Health check
- `GET /health-check` — Service health status

### Customer endpoints
- `POST /customer/validate-create` — Validate customer data before creation
- `POST /customer/create` — Create a new customer

### Address endpoints
- `POST /address/create` — Create a new address
- `GET /address/search/{id}` — Get address by ID

### Phone endpoints
Phone management is integrated into customer operations with support for:
- Mobile phone validation (9 digits)
- Landline validation (8 digits)
- Brazilian area code (DDD) validation

---

## Architecture

The codebase follows a Hexagonal architecture (Ports & Adapters) with clear separation of concerns:

### Project structure
```
petshop-api/
├── adapter/
│   ├── input/
│   │   ├── http/          # HTTP handlers and routing
│   │   └── message/       # Kafka consumers
│   └── output/
│       ├── cache/         # Redis implementations
│       └── database/      # PostgreSQL repositories
├── application/
│   ├── domain/            # Domain models and context
│   ├── port/
│   │   ├── input/         # Use case interfaces
│   │   └── output/        # Repository interfaces
│   ├── service/           # Business logic implementation
│   └── utils/             # Validation utilities (CPF, CNPJ, etc.)
├── cmd/
│   └── petshop-api/       # Application entry point
├── configuration/
│   ├── db/                # SQL scripts and migrations
│   ├── environment/       # Environment configuration
│   └── repository/        # Database connection setup
└── diagrams/              # Architecture diagrams
```

### Design principles
- **Ports & Adapters**: Business logic is isolated from infrastructure concerns
- **Dependency Inversion**: Interfaces define contracts; implementations are injected
- **Testability**: Comprehensive mocks and test coverage
- **Caching Strategy**: Non-fatal cache failures; cache is a performance optimization
- **Error Handling**: Structured error handling with proper logging

---

## Validation features

The API includes robust validation for Brazilian-specific data:

### Document validation
- **CPF** (Cadastro de Pessoas Físicas) — Individual taxpayer ID
  - Length validation (11 digits)
  - Check digit verification
  - Invalid patterns detection (e.g., all same digits)

- **CNPJ** (Cadastro Nacional da Pessoa Jurídica) — Company taxpayer ID
  - Length validation (14 digits)
  - Check digit verification
  - Invalid patterns detection

### Phone validation
- **Area codes (DDD)** — All valid Brazilian area codes supported
- **Mobile phones** — 9-digit validation
- **Landline phones** — 8-digit validation

### Address validation
- **Required fields** — Street, Number, Neighborhood, ZipCode, City, State, Country
- **State field** — Must be exactly 2 characters (Brazilian state codes: RJ, SP, MG, etc.)
- **Comprehensive error reporting** — Returns all validation failures in a single error message

---

## Database

### Schema
The database includes the following main tables:

**petshop_api schema**
- `address` — Address information with Brazilian format validation
- `customer` — Customer data with CPF/CNPJ validation
- `phone` — Phone contacts with DDD and number type
- `contract` — Contract information for legal entities

**petshop_auth schema**
- Authentication and authorization tables (managed by gateway)

**petshop_gateway schema**
- API gateway configuration and routing

### Initialize database
SQL initialization scripts are located in `configuration/db/` and are automatically executed when starting the Docker environment. These scripts include:
- Schema creation (tables, constraints, indexes)
- Seed data for testing and development

---

## Testing

The project includes comprehensive test coverage:

- Unit tests for services (`application/service/*_test.go`)
- Unit tests for utilities (`application/utils/utils_test.go`)
- Integration tests for HTTP handlers (`adapter/input/http/handler/*_test.go`)
- Mock implementations for all repository interfaces

### Test naming convention
Tests follow the pattern: `TestFunction_Scenario_ExpectedBehavior`

Examples:
```go
TestAddressService_Create/WithValidAddress_SavesSuccessfully
TestAddressService_Create/WithInvalidAddress_ValidationFails
TestAddressService_ValidateAddress/WithInvalidStateLength_ReturnsError
TestAddressService_GetById/WithValidID_ReturnsAddress
```

### Run specific tests
```bash
go test ./application/service -v
go test ./adapter/input/http/handler -v
go test ./application/utils -v
```

### Coverage report
```bash
make test-cover
```
This generates an HTML coverage report (excluding mocks) and opens it in your browser.

---

## Logging

The application uses structured logging with zap:

- JSON format for production
- ISO 8601 timestamps
- Contextual fields (customer_id, address_id, etc.)
- Log levels: DEBUG, INFO, WARN, ERROR
- Caller information for error tracing

---

## Contributing

When contributing to this repository:

1. Follow Go best practices and idioms
2. Write tests for new features (follow the naming convention)
3. Use the existing code structure (Hexagonal architecture)
4. Add proper error handling and logging
5. Update documentation as needed
6. Run `go fmt` and linters before committing
7. Ensure all tests pass with `go test ./...`
8. Consider non-breaking changes for validation rules

---

## License

This project is available under the MIT License — see the `LICENSE` file for details.

---

