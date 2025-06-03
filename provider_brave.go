package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func fetchBrave(query string) ([]SearchResult, error) {
	braveURL := fmt.Sprintf("https://search.brave.com/search?q=%s&offset=0", url.QueryEscape(query))

	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", braveURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Brave request: %w", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:139.0) Gecko/20100101 Firefox/139.0")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "https://search.brave.com/")
	req.Header.Add("Cookie", "safe_search=off")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Brave: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Brave search returned status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Brave HTML response: %w", err)
	}

	var results []SearchResult

	doc.Find("#results [data-pos]").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".h .url .title").Text()
		url, _ := s.Find("a").Attr("href")
		description := s.Find(".snippet-content").Text()

		if title != "" && url != "" {
			results = append(results, SearchResult{
				Title:       title,
				URL:         url,
				Description: description,
				Engine:      "brave",
				Category:    "Web",
			})
		}
	})

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found from brave")
	}

	return results, nil
}
