package interfaces

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/t-yamakoshi/go-mcp-client/pkg/domain/entity"
	"github.com/t-yamakoshi/go-mcp-client/pkg/usecase"
)

// CLIHandler defines the interface for command line interface operations
type CLIHandler interface {
	Run() error
}

// cliHandler implements the CLIHandler interface
type cliHandler struct {
	mcpUseCase    usecase.MCPUseCase
	configUseCase usecase.ConfigUseCase
}

// NewCLIHandler creates a new CLI handler
func NewCLIHandler(mcpUseCase usecase.MCPUseCase, configUseCase usecase.ConfigUseCase) CLIHandler {
	return &cliHandler{
		mcpUseCase:    mcpUseCase,
		configUseCase: configUseCase,
	}
}

// Run starts the CLI application
func (h *cliHandler) Run() error {
	// Parse command line flags
	configFile := flag.String("config", "config.json", "Path to configuration file")
	serverURL := flag.String("server", "", "MCP server URL (overrides config file)")
	flag.Parse()

	// Load configuration
	config, err := h.configUseCase.LoadConfiguration(context.Background(), *configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override server URL if provided via command line
	if *serverURL != "" {
		config.ServerURL = *serverURL
	}

	// Set up context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, closing connection...")
		cancel()
	}()

	// Establish connection
	log.Printf("Connecting to MCP server at %s", config.ServerURL)
	if err := h.mcpUseCase.EstablishConnection(ctx, config.ServerURL); err != nil {
		return fmt.Errorf("failed to connect to MCP server: %w", err)
	}
	defer h.mcpUseCase.CloseConnection(ctx)

	// Initialize the protocol
	log.Println("Initializing MCP protocol...")
	initResp, err := h.mcpUseCase.InitializeProtocol(ctx, config.ClientInfo)
	if err != nil {
		return fmt.Errorf("failed to initialize MCP protocol: %w", err)
	}

	log.Printf("Successfully connected to MCP server: %s v%s",
		initResp.ServerInfo.Name, initResp.ServerInfo.Version)

	// Register message handlers
	h.registerHandlers()

	// Keep the connection alive
	log.Println("MCP client is running. Press Ctrl+C to exit.")
	<-ctx.Done()
	log.Println("MCP client shutting down...")

	return nil
}

// registerHandlers registers message handlers for the MCP client
func (h *cliHandler) registerHandlers() {
	// Register handler for tools/list
	h.mcpUseCase.RegisterHandler("tools/list", func(msg *entity.Message) error {
		log.Println("Received tools/list request")
		// Here you would implement the actual tools listing logic
		return nil
	})

	// Register handler for tools/call
	h.mcpUseCase.RegisterHandler("tools/call", func(msg *entity.Message) error {
		log.Println("Received tools/call request")
		// Here you would implement the actual tool calling logic
		return nil
	})

	// Register handler for ping
	h.mcpUseCase.RegisterHandler("ping", func(msg *entity.Message) error {
		log.Println("Received ping, sending pong")
		// Send pong response
		pongMsg := &entity.Message{
			ID:     msg.ID,
			Method: "pong",
		}
		return h.mcpUseCase.SendOutgoingMessage(context.Background(), pongMsg)
	})
}

// HTTPHandler defines the interface for HTTP interface operations (for future use)
type HTTPHandler interface {
	StartServer(ctx context.Context, port string) error
}

// httpHandler implements the HTTPHandler interface
type httpHandler struct {
	mcpUseCase    usecase.MCPUseCase
	configUseCase usecase.ConfigUseCase
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(mcpUseCase usecase.MCPUseCase, configUseCase usecase.ConfigUseCase) HTTPHandler {
	return &httpHandler{
		mcpUseCase:    mcpUseCase,
		configUseCase: configUseCase,
	}
}

// StartServer starts the HTTP server (placeholder for future implementation)
func (h *httpHandler) StartServer(ctx context.Context, port string) error {
	// TODO: Implement HTTP server
	log.Printf("HTTP server would start on port %s", port)
	return nil
}
