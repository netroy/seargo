package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type SearchArgs struct {
	Query   string `json:"query"`
	Engines string `json:"engines"`
}

func mcpHandler(ctx context.Context, request mcp.CallToolRequest, args SearchArgs) (*mcp.CallToolResult, error) {
	if args.Engines == "" {
		args.Engines = "brave"
	}
	selectedEngines := strings.Split(args.Engines, ",")

	results := runQuery(args.Query, selectedEngines)

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

func NewMCPServer() *server.MCPServer {
	mcpServer := server.NewMCPServer("seargo", "0.1.0")

	tool := mcp.NewTool(
		"search",
		mcp.WithDescription("Performs a search query across multiple engines."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The search query string"),
		),
		mcp.WithString("engines",
			mcp.Description("Comma-separated list of engines to query (e.g., 'duckduckgo,brave,wikipedia'). Defaults to 'brave'."),
		),
	)

	mcpServer.AddTool(tool, mcp.NewTypedToolHandler(mcpHandler))

	return mcpServer
}
