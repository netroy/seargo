package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
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

	var (
		allResults []SearchResult
		mu         sync.Mutex
		wg         sync.WaitGroup
	)

	_, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	type searchProvider func(query string) ([]SearchResult, error)

	engineMap := map[string]searchProvider{
		"duckduckgo": fetchDuckDuckGo,
		"brave":      fetchBrave,
		"wikipedia":  fetchWikipedia,
	}

	for _, engineName := range selectedEngines {
		if provider, ok := engineMap[engineName]; ok {
			wg.Add(1)
			go func(p searchProvider) {
				defer wg.Done()
				results, err := p(query)
				if err != nil {
					log.Printf("Search provider error (%s): %v", engineName, err)
					return
				}
				mu.Lock()
				allResults = append(allResults, results...)
				mu.Unlock()
			}(provider)
		} else {
			log.Printf("Unknown engine: %s", engineName)
		}
	}

	wg.Wait() // Wait for all goroutines to finish

	if format == "json" {
		return c.JSON(allResults)
	}

	return c.Render("results", fiber.Map{
		"Query":     query,
		"Results":   allResults,
		"FetchTime": time.Since(startTime).Seconds(),
	})
}

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./static")

	app.Get("/search", searchHandler)

	log.Fatal(app.Listen(":8080"))
}
