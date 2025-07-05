package cli

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/t-yamakoshi/go-mcp-client/pkg/interfaces/message"
	"github.com/t-yamakoshi/go-mcp-client/pkg/usecase"
)

var _ IFCLIHandler = (*CliHandler)(nil)

// IFCLIHandler defines the interface for command line interface operations
type IFCLIHandler interface {
	Run() error
}

// CliHandler implements the IFCLIHandler interface
type CliHandler struct {
	mcpUsecase    usecase.IFMCPUsecase
	configUsecase usecase.IFConfigUsecase
	msgHandler    message.IFMessageHandler
}

// NewCLIHandler creates a new CLI handler
func NewCLIHandler(mcpUsecase *usecase.MCPUsecase, configUsecase *usecase.ConfigUsecase, msgHandler *message.MessageHandler) *CliHandler {
	return &CliHandler{
		mcpUsecase:    mcpUsecase,
		configUsecase: configUsecase,
		msgHandler:    msgHandler,
	}
}

// Run starts the CLI application
func (h *CliHandler) Run() error {
	// Parse command line flags
	configFile := flag.String("config", "config.json", "Path to configuration file")
	serverURL := flag.String("server", "", "MCP server URL (overrides config file)")
	flag.Parse()

	// Load configuration
	config, err := h.configUsecase.LoadConfiguration(context.Background(), *configFile)
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
	if err := h.mcpUsecase.EstablishConnection(ctx, config.ServerURL); err != nil {
		return fmt.Errorf("failed to connect to MCP server: %w", err)
	}
	defer h.mcpUsecase.CloseConnection(ctx)

	// Initialize the protocol
	log.Println("Initializing MCP protocol...")
	initResp, err := h.mcpUsecase.InitializeProtocol(ctx, config.ClientInfo)
	if err != nil {
		return fmt.Errorf("failed to initialize MCP protocol: %w", err)
	}

	log.Printf("Successfully connected to MCP server: %s v%s",
		initResp.ServerInfo.Name, initResp.ServerInfo.Version)

	// Register message handlers
	h.msgHandler.RegisterHandlers(&h.mcpUsecase)

	// Keep the connection alive
	log.Println("MCP client is running. Press Ctrl+C to exit.")
	<-ctx.Done()
	log.Println("MCP client shutting down...")

	return nil
}
