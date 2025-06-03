package main

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ProviderConfig struct {
	BaseURL             string
	ResultSelector      string
	TitleSelector       string
	URLSelector         string
	DescriptionSelector string
}

type SearchResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

var providersMap = map[string]ProviderConfig{
	"bing": {
		BaseURL:             "https://www.bing.com/search?q=%s",
		ResultSelector:      "li.b_algo",
		TitleSelector:       "h2 a",
		URLSelector:         "h2 a",
		DescriptionSelector: "p",
	},
	"brave": {
		BaseURL:             "https://search.brave.com/search?q=%s&offset=0",
		ResultSelector:      "#results [data-pos]",
		TitleSelector:       ".h .url .title",
		URLSelector:         "a",
		DescriptionSelector: ".snippet-content",
	},
	"duckduckgo": {
		BaseURL:             "https://html.duckduckgo.com/html/?q=%s",
		ResultSelector:      ".results .result",
		TitleSelector:       ".result__title a",
		URLSelector:         ".result__title a",
		DescriptionSelector: ".result__snippet",
	},
	"startpage": {
		BaseURL:             "https://www.startpage.com/do/search?query=%s",
		ResultSelector:      "div.w-gl .result",
		TitleSelector:       ".result-title .wgl-title",
		URLSelector:         "a.result-link",
		DescriptionSelector: "p.description",
	},
	"wikipedia": {
		BaseURL:             "https://en.wikipedia.org/w/index.php?limit=10&offset=0&ns0=1&search=%s",
		ResultSelector:      ".mw-search-result",
		TitleSelector:       ".mw-search-result-heading a",
		URLSelector:         ".mw-search-result-heading a",
		DescriptionSelector: ".searchresult",
	},
}

func runQuery(query string, providerName string) []SearchResult {
	results, err := scrapeWebSearch(query, providerName)
	if err != nil {
		log.Printf("Search provider error (%s): %v", providerName, err)
		return []SearchResult{}
	}
	return results
}

func scrapeWebSearch(query string, providerName string) ([]SearchResult, error) {
	providerConfig, _ := providersMap[providerName]

	searchURL := fmt.Sprintf(providerConfig.BaseURL, url.QueryEscape(query))

	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", providerName, err)
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", config.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from %s: %w", providerName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s search returned status %d", providerName, resp.StatusCode)
	}

	var reader io.Reader
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	} else if resp.Header.Get("Content-Encoding") == "deflate" {
		reader = flate.NewReader(resp.Body)
	} else {
		reader = resp.Body
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s HTML response: %w", providerName, err)
	}

	htmlContent, err := doc.Html()
	if err != nil {
		return nil, fmt.Errorf("failed to get HTML content from %s: %w", providerName, err)
	}

	results, err := parseWebSearch(htmlContent, providerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s search results: %w", providerName, err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found from %s", providerName)
	}

	return results, nil
}

func parseWebSearch(htmlContent string, providerConfig ProviderConfig) ([]SearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to create goquery document from HTML content: %w", err)
	}

	var results []SearchResult

	doc.Find(providerConfig.ResultSelector).Each(func(i int, s *goquery.Selection) {
		title := s.Find(providerConfig.TitleSelector).Text()
		url, _ := s.Find(providerConfig.URLSelector).Attr("href")
		description := s.Find(providerConfig.DescriptionSelector).Text()

		if title != "" && url != "" {
			results = append(results, SearchResult{
				Title:       title,
				URL:         url,
				Description: description,
			})
		}
	})

	if len(results) == 0 {
		return nil, nil
	}
	return results, nil
}
