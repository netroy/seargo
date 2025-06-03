package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"slices"
	"time"
)

type WebSearchHandler struct{}

type WebSearchData struct {
	Provider  string
	Query     string
	Results   []SearchResult
	FetchTime float64
}

func (h WebSearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	query := queryParams.Get("q")
	if query == "" {
		w.WriteHeader(400)
		w.Write([]byte("Missing query parameter 'q'"))
		return
	}

	provider := queryParams.Get("provider")
	if provider == "" {
		provider = config.DefaultProvider
	}
	if !slices.Contains(config.EnabledProviders, provider) {
		w.WriteHeader(400)
		w.Write([]byte("Invalid search provider"))
	}

	format := queryParams.Get("format")
	startTime := time.Now()
	results := runQuery(query, provider)

	if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
		return
	}

	cwd, _ := os.Getwd()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, _ := template.ParseFiles(fmt.Sprintf("%s/templates/results.html", cwd))
	t.Execute(w, WebSearchData{
		Provider:  provider,
		Query:     query,
		Results:   results,
		FetchTime: time.Since(startTime).Seconds(),
	})
}

func setupWebServer(mux *http.ServeMux) {
	mux.Handle("/search", WebSearchHandler{})

	cwd, _ := os.Getwd()
	static := http.FileServer(http.Dir(fmt.Sprintf("%s/static", cwd)))
	mux.Handle("/", static)
	mux.Handle("/favicon.ico", static)
	mux.Handle("/static/", http.StripPrefix("/static/", static))
}
