package http

import (
	"context"
	"log"

	"github.com/t-yamakoshi/go-mcp-client/pkg/usecase"
)

var _ IFHTTPHandler = (*HTTPHandler)(nil)

type IFHTTPHandler interface {
	StartServer(ctx context.Context, port string) error
}

type HTTPHandler struct {
	mcpUsecase    usecase.IFMCPUsecase
	configUsecase usecase.IFConfigUsecase
}

func NewHTTPHandler(mcpUsecase *usecase.MCPUsecase, configUsecase *usecase.ConfigUsecase) *HTTPHandler {
	return &HTTPHandler{
		mcpUsecase:    mcpUsecase,
		configUsecase: configUsecase,
	}
}

func (h *HTTPHandler) StartServer(ctx context.Context, port string) error {
	log.Printf("HTTP server would start on port %s", port)
	return nil
}
