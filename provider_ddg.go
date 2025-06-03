package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

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
		return nil, fmt.Errorf("no results found on duckduckgo")
	}

	return results, nil
}
