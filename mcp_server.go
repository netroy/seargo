package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type SearchArgs struct {
	Query    string `json:"query"`
	Provider string `json:"provider"`
}

func mcpHandler(ctx context.Context, request mcp.CallToolRequest, args SearchArgs) (*mcp.CallToolResult, error) {
	provider := args.Provider
	if provider == "" {
		if len(config.EnabledProviders) > 0 {
			provider = config.EnabledProviders[0]
		} else {
			return nil, fmt.Errorf("no provider specified and no enabled providers configured")
		}
	}

	results := runQuery(args.Query, provider)

	jsonResults, err := json.Marshal(results)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search results: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(jsonResults),
			},
		},
	}, nil
}

func setupMCPServer(httpServer *http.Server, mux *http.ServeMux) {
	mcpServer := server.NewMCPServer("seargo", "0.1.0")

	tool := mcp.NewTool(
		"search",
		mcp.WithDescription("Performs a search query across multiple providers."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The search query string"),
		),
		mcp.WithString("provider",
			mcp.Description("The single provider to query (e.g., 'duckduckgo', 'brave', 'bing). Defaults to Bing."),
		),
	)

	mcpServer.AddTool(tool, mcp.NewTypedToolHandler(mcpHandler))

	sseServer := server.NewSSEServer(
		mcpServer,
		server.WithBaseURL(config.BaseURL),
		server.WithHTTPServer(httpServer),
	)
	mux.Handle("/sse", sseServer)
	mux.Handle("/message", sseServer)

	streamableHttpServer := server.NewStreamableHTTPServer(
		mcpServer,
		server.WithStreamableHTTPServer(httpServer),
		server.WithStateLess(true),
	)
	mux.Handle("/mcp", streamableHttpServer)
}
