package main

import (
	"log"
	"sync"
)

type searchProvider func(query string) ([]SearchResult, error)

func runQuery(query string, selectedEngines []string) []SearchResult {
	var (
		allResults []SearchResult
		mu         sync.Mutex
		wg         sync.WaitGroup
	)

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
	return allResults
}
