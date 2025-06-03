package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./static")

	app.Get("/search", func(c *fiber.Ctx) error {
		query := c.Query("q")
		if query == "" {
			return c.Status(400).SendString("Missing query parameter 'q'")
		}

		format := c.Query("format", "html")

		startTime := time.Now()

		var (
			allResults []SearchResult
			mu         sync.Mutex
			wg         sync.WaitGroup
		)

		_, cancel := context.WithTimeout(c.Context(), 10*time.Second)
		defer cancel()

		// Fetch from DuckDuckGo
		wg.Add(1)
		go func() {
			defer wg.Done()
			ddgResults, err := fetchDuckDuckGo(query)
			if err != nil {
				log.Printf("DuckDuckGo search error: %v", err)
				return
			}
			mu.Lock()
			allResults = append(allResults, ddgResults...)
			mu.Unlock()
		}()

		// Fetch from Brave
		wg.Add(1)
		go func() {
			defer wg.Done()
			braveResults, err := fetchBrave(query)
			if err != nil {
				log.Printf("Brave search error: %v", err)
				return
			}
			mu.Lock()
			allResults = append(allResults, braveResults...)
			mu.Unlock()
		}()

		wg.Wait() // Wait for all goroutines to finish

		if format == "json" {
			return c.JSON(allResults)
		}

		return c.Render("results", fiber.Map{
			"Query":     query,
			"Results":   allResults,
			"FetchTime": time.Since(startTime).Seconds(),
		})
	})

	log.Fatal(app.Listen(":8080"))
}
