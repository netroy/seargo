package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type SearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
	Engine  string `json:"engine"`
}

func fetchDuckDuckGo(query string) ([]SearchResult, error) {
	apiURL := fmt.Sprintf(
		"https://api.duckduckgo.com/?q=%s&format=json&no_html=1&no_redirect=1",
		url.QueryEscape(query),
	)

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from DuckDuckGo: %w", err)
	}
	defer resp.Body.Close()

	var ddgResponse struct {
		RelatedTopics []struct {
			Name     string `json:"Name"`
			Text     string `json:"Text"`
			FirstURL string `json:"FirstURL"`
		} `json:"RelatedTopics"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ddgResponse); err != nil {
		return nil, fmt.Errorf("failed to decode DuckDuckGo response: %w", err)
	}

	var results []SearchResult
	for _, topic := range ddgResponse.RelatedTopics {
		if topic.FirstURL != "" {
			results = append(results, SearchResult{
				Title:   topic.Name,
				URL:     topic.FirstURL,
				Content: topic.Text,
				Engine:  "duckduckgo",
			})
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found")
	}

	return results, nil
}

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/search", func(c *fiber.Ctx) error {
		query := c.Query("q")
		if query == "" {
			return c.Status(400).SendString("Missing query parameter 'q'")
		}

		format := c.Query("format", "html")

		results, err := fetchDuckDuckGo(query)
		if err != nil {
			log.Printf("Search error: %v", err)
			return c.Status(500).SendString("Search failed")
		}

		if format == "json" {
			return c.JSON(results)
		}

		return c.Render("results", fiber.Map{
			"Query":   query,
			"Results": results,
		})
	})

	log.Fatal(app.Listen(":8080"))
}
