package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func fetchWikipedia(query string) ([]SearchResult, error) {
	host := "https://en.wikipedia.org"

	queryParams := url.Values{}
	queryParams.Add("limit", "20")
	queryParams.Add("offset", "0")
	queryParams.Add("profile", "default")
	queryParams.Add("search", query)
	queryParams.Add("title", "Special:Search")
	queryParams.Add("ns0", "1")

	wikipediaURL := fmt.Sprintf("%s/w/index.php?%s", host, queryParams.Encode())

	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", wikipediaURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create wikipedia request: %w", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:139.0) Gecko/20100101 Firefox/139.0")
	req.Header.Add("Referer", host)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from wikipedia: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wikipedia search returned status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse wikipedia HTML response: %w", err)
	}

	var results []SearchResult

	doc.Find(".mw-search-results li.mw-search-result").Each(func(i int, s *goquery.Selection) {
		titleSelection := s.Find(".mw-search-result-heading a")
		title := titleSelection.Text()
		relativeURL, _ := titleSelection.Attr("href")
		description := s.Find(".searchresult").Text()

		if title != "" && relativeURL != "" {
			fullURL := fmt.Sprintf("%s%s", host, relativeURL)
			results = append(results, SearchResult{
				Title:       title,
				URL:         fullURL,
				Description: description,
				Engine:      "wikipedia",
				Category:    "Knowledge",
			})
		}
	})

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found from wikipedia")
	}

	return results, nil
}
