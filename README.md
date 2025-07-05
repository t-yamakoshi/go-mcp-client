# Go MCP Client

A Go implementation of the Model Context Protocol (MCP) client built with Clean Architecture principles.

## Overview

This project implements a client for the Model Context Protocol, which allows AI assistants to connect to external data sources and tools through a standardized interface. The application is structured following Clean Architecture principles to ensure maintainability, testability, and separation of concerns.

## Architecture

This project follows Clean Architecture principles with the following layers:

### 1. Domain Layer (`pkg/domain/`)
- **Entities**: Core business objects (Message, Tool, Connection, etc.)
- **Repository Interfaces**: Abstract interfaces for data access
- **Service Interfaces**: Abstract interfaces for business logic

### 2. Use Case Layer (`pkg/usecase/`)
- **Business Logic**: Application-specific business rules
- **Orchestration**: Coordinates between different domain services
- **Input/Output Ports**: Defines how the application interacts with external systems

### 3. Interface Layer (`pkg/interfaces/`)
- **CLI Handler**: Command-line interface implementation
- **HTTP Handler**: HTTP interface implementation (future)
- **Controllers**: Handle user input and format output

### 4. Infrastructure Layer (`pkg/infrastructure/`)
- **MCP Repository**: WebSocket implementation for MCP protocol
- **Config Repository**: File-based configuration storage
- **External Services**: Database, external APIs, etc.

## Features

- WebSocket-based communication with MCP servers
- Support for MCP protocol version 2024-11-05
- Configurable client settings
- Graceful shutdown handling
- Extensible message handler system
- Clean Architecture design
- Dependency injection
- Separation of concerns

## Installation

```bash
git clone <repository-url>
cd go-mcp-client
go mod tidy
```

## Usage

### Basic Usage

```bash
# Run with default configuration
go run cmd/main.go

# Run with custom server URL
go run cmd/main.go -server ws://localhost:3000

# Run with custom config file
go run cmd/main.go -config my-config.json
```

### Configuration

The client can be configured using a JSON configuration file. If no configuration file is provided, a default configuration will be created.

Example configuration (`config.json`):

```json
{
  "server_url": "ws://localhost:3000",
  "client_info": {
    "name": "go-mcp-client",
    "version": "1.0.0"
  },
  "log_level": "info"
}
```

### Command Line Flags

- `-config`: Path to configuration file (default: `config.json`)
- `-server`: MCP server URL (overrides config file)

## Project Structure

```
go-mcp-client/
├── cmd/
│   └── main.go                    # Main entry point with DI setup
├── pkg/
│   ├── domain/                    # Domain Layer
│   │   ├── entity.go              # Core business entities
│   │   ├── repository.go          # Repository interfaces
│   │   └── service.go             # Service interfaces
│   ├── usecase/                   # Use Case Layer
│   │   ├── mcp_usecase.go         # MCP business logic
│   │   └── config_usecase.go      # Configuration business logic
│   ├── interfaces/                # Interface Layer
│   │   └── handler.go             # CLI and HTTP handlers
│   └── infrastructure/            # Infrastructure Layer
│       ├── mcp_repository.go      # MCP WebSocket implementation
│       └── config_repository.go   # File-based config storage
├── go.mod
└── README.md
```

## Dependency Flow

The dependency flow follows Clean Architecture principles:

```
Interfaces → Use Cases → Domain ← Infrastructure
     ↓           ↓         ↑           ↑
   (Input)   (Business)  (Core)    (External)
```

- **Interfaces** depend on **Use Cases**
- **Use Cases** depend on **Domain** interfaces
- **Infrastructure** implements **Domain** interfaces
- **Domain** has no dependencies on other layers

## Development

### Building

```bash
go build -o mcp-client cmd/main.go
```

### Running Tests

```bash
go test ./...
```

### Adding New Features

1. **Domain Layer**: Define entities and interfaces
2. **Use Case Layer**: Implement business logic
3. **Interface Layer**: Add user interface handlers
4. **Infrastructure Layer**: Implement external integrations

## MCP Protocol Support

This client supports the following MCP protocol features:

- Connection establishment and initialization
- Message handling with custom handlers
- Tool listing and calling (framework ready)
- Ping/pong heartbeat mechanism

## Benefits of Clean Architecture

- **Testability**: Each layer can be tested independently
- **Maintainability**: Clear separation of concerns
- **Flexibility**: Easy to swap implementations
- **Scalability**: Well-defined boundaries between components
- **Independence**: Business logic is independent of external frameworks

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes following Clean Architecture principles
4. Add tests for new functionality
5. Submit a pull request

## License

[Add your license here] 
