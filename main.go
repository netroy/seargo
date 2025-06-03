package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/mark3labs/mcp-go/server"
)

func parseEngines(enginesParam string) []string {
	if enginesParam == "" {
		return []string{}
	}
	return strings.Split(enginesParam, ",")
}

func searchHandler(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(400).SendString("Missing query parameter 'q'")
	}

	enginesParam := c.Query("engines", "brave") // Default to "brave"
	selectedEngines := parseEngines(enginesParam)

	format := c.Query("format", "html")

	startTime := time.Now()

	_, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	results := runQuery(query, selectedEngines)

	if format == "json" {
		return c.JSON(results)
	}

	return c.Render("results", fiber.Map{
		"Query":     query,
		"Results":   results,
		"FetchTime": time.Since(startTime).Seconds(),
	})
}

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(logger.New())

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: adaptor.FiberApp(app),
	}
	mcpServer := NewMCPServer()

	// server.NewStreamableHTTPServer(
	// 	mcpServer,
	// 	server.WithStreamableHTTPServer(httpServer),
	// 	server.WithEndpointPath("mcp"),
	// )
	sseServer := server.NewSSEServer(
		mcpServer,
		server.WithHTTPServer(httpServer),
		server.WithStaticBasePath("/mcp"),
		server.WithBaseURL("http://localhost:8080"),
	)
	app.All("/mcp/sse", adaptor.HTTPHandler(sseServer.SSEHandler()))
	app.All("/mcp/message", adaptor.HTTPHandler(sseServer.MessageHandler()))

	cwd, _ := os.Getwd()
	app.Get("/search", searchHandler)
	app.Static("/", fmt.Sprintf("%s/static", cwd))

	log.Printf("Server listening on %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
