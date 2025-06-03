package main

import (
	"reflect"
	"testing"
)

func TestParseWebSearch(t *testing.T) {
	tests := []struct {
		name           string
		htmlContent    string
		providerConfig ProviderConfig
		want           []SearchResult
		wantErr        bool
	}{
		{
			name: "Successful parsing with multiple results",
			htmlContent: `
				<html>
				<body>
					<div class="result-item">
						<h3 class="result-title"><a href="http://example.com/page1">Title 1</a></h3>
						<p class="result-description">Description 1</p>
					</div>
					<div class="result-item">
						<h3 class="result-title"><a href="http://example.com/page2">Title 2</a></h3>
						<p class="result-description">Description 2</p>
					</div>
				</body>
				</html>
			`,
			providerConfig: ProviderConfig{
				ResultSelector:      ".result-item",
				TitleSelector:       ".result-title a",
				URLSelector:         ".result-title a",
				DescriptionSelector: ".result-description",
			},
			want: []SearchResult{
				{Title: "Title 1", URL: "http://example.com/page1", Description: "Description 1"},
				{Title: "Title 2", URL: "http://example.com/page2", Description: "Description 2"},
			},
			wantErr: false,
		},
		{
			name: "No results found",
			htmlContent: `
				<html>
				<body>
					<div>No results here</div>
				</body>
				</html>
			`,
			providerConfig: ProviderConfig{
				ResultSelector:      ".non-existent-result",
				TitleSelector:       "a",
				URLSelector:         "a",
				DescriptionSelector: "p",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Partial result (missing URL)",
			htmlContent: `
				<html>
				<body>
					<div class="result-item">
						<h3 class="result-title"><a>Title 1</a></h3>
						<p class="result-description">Description 1</p>
					</div>
				</body>
				</html>
			`,
			providerConfig: ProviderConfig{
				ResultSelector:      ".result-item",
				TitleSelector:       ".result-title a",
				URLSelector:         ".result-title a",
				DescriptionSelector: ".result-description",
			},
			want:    nil, // Should not append if URL is empty
			wantErr: false,
		},
		{
			name: "Partial result (missing Title)",
			htmlContent: `
				<html>
				<body>
					<div class="result-item">
						<h3 class="result-title"><a href="http://example.com/page1"></a></h3>
						<p class="result-description">Description 1</p>
					</div>
				</body>
				</html>
			`,
			providerConfig: ProviderConfig{
				ResultSelector:      ".result-item",
				TitleSelector:       ".result-title a",
				URLSelector:         ".result-title a",
				DescriptionSelector: ".result-description",
			},
			want:    nil, // Should not append if Title is empty
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseWebSearch(tt.htmlContent, tt.providerConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseWebSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseWebSearch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
