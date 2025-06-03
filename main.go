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
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	Engine      string `json:"engine"`
}

type DDGResponse struct {
	Abstract       string `json:"Abstract"`
	AbstractText   string `json:"AbstractText"`
	AbstractURL    string `json:"AbstractURL"`
	AbstractSource string `json:"AbstractSource"`
	Image          string `json:"Image"`
	Heading        string `json:"Heading"`
	Answer         string `json:"Answer"`
	Redirect       string `json:"Redirect"`
	RelatedTopics  []struct {
		Result string `json:"Result"`
		Icon   struct {
			URL string `json:"URL"`
		} `json:"Icon"`
		FirstURL string `json:"FirstURL"`
		Text     string `json:"Text"`
	} `json:"RelatedTopics"`
	Results []struct {
		Result string `json:"Result"`
		Icon   struct {
			URL    string `json:"URL"`
			Height int    `json:"Height"`
			Width  int    `json:"Width"`
		} `json:"Icon"`
		FirstURL string `json:"FirstURL"`
		Text     string `json:"Text"`
	} `json:"Results"`
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

	var ddgResponse DDGResponse
	if err := json.NewDecoder(resp.Body).Decode(&ddgResponse); err != nil {
		return nil, fmt.Errorf("failed to decode DuckDuckGo response: %w", err)
	}

	var results []SearchResult

	// Parse abstract/infobox result
	if ddgResponse.AbstractURL != "" {
		results = append(results, SearchResult{
			Title:       ddgResponse.Heading,
			URL:         ddgResponse.AbstractURL,
			Description: ddgResponse.AbstractText,
			Icon:        ddgResponse.Image,
			Category:    "Infobox",
			Engine:      "duckduckgo",
		})
	}

	// Parse related topics
	for _, topic := range ddgResponse.RelatedTopics {
		if topic.FirstURL != "" {
			results = append(results, SearchResult{
				Title:       topic.Text,
				URL:         topic.FirstURL,
				Description: topic.Text,
				Icon:        topic.Icon.URL,
				Category:    "Related",
				Engine:      "duckduckgo",
			})
		}
	}

	// Parse regular results
	for _, result := range ddgResponse.Results {
		if result.FirstURL != "" {
			results = append(results, SearchResult{
				Title:       result.Text,
				URL:         result.FirstURL,
				Description: result.Text,
				Icon:        result.Icon.URL,
				Category:    "Web",
				Engine:      "duckduckgo",
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

		startTime := time.Now()
		results, err := fetchDuckDuckGo(query)
		if err != nil {
			log.Printf("Search error: %v", err)
			return c.Status(500).SendString("Search failed")
		}

		if format == "json" {
			return c.JSON(results)
		}

		return c.Render("results", fiber.Map{
			"Query":     query,
			"Results":   results,
			"FetchTime": time.Now().Sub(startTime).Seconds(),
		})
	})

	log.Fatal(app.Listen(":8080"))
}
