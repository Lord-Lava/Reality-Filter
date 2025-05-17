# Reality Filter

A backend system for detecting and analyzing potentially fake, misleading, or biased news articles using hexagonal architecture.

## Project Structure

```
├── cmd/                    # Application entry points
│   └── server/            # Main server application
├── internal/              # Private application code
│   ├── core/             # Domain layer (business logic)
│   │   ├── domain/       # Domain models and business rules
│   │   └── ports/        # Interfaces (ports)
│   ├── adapters/         # Adapters implementation
│   │   ├── primary/      # Inbound adapters (REST, gRPC, Kafka consumers)
│   │   └── secondary/    # Outbound adapters (Redis, PostgreSQL, External APIs)
│   └── application/      # Application services
├── pkg/                  # Public library code
│   ├── config/          # Configuration
│   └── logger/          # Logging utilities
└── api/                 # API contracts and proto files
    ├── rest/            # REST API specifications
    └── proto/           # Protocol buffer definitions
```

## Hexagonal Architecture Components

1. **Domain Layer** (`internal/core/domain/`)
   - Article entities and value objects
   - Scoring and analysis logic
   - Business rules and validation

2. **Ports** (`internal/core/ports/`)
   - Primary (Inbound) Ports: Service interfaces for article analysis
   - Secondary (Outbound) Ports: Repository and external service interfaces

3. **Adapters** (`internal/adapters/`)
   - Primary Adapters: HTTP handlers, gRPC servers, Kafka consumers
   - Secondary Adapters: Database implementations, cache, external API clients

4. **Application Services** (`internal/application/`)
   - Use case implementations
   - Orchestration of domain logic
   - Transaction management

## Getting Started

1. Install dependencies:
```bash
go mod download
```

2. Run the server:
```bash
go run cmd/server/main.go
```

## Development

- Follow Go best practices and project layout conventions
- Use dependency injection for adapters
- Keep the domain layer pure and independent
- Write tests for each layer independently 