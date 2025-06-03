package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	BaseURL          string   `env:"BASE_URL"`
	Port             int      `env:"PORT, default=8080"`
	EnabledProviders []string `env:"ENABLED_PROVIDERS, default=bing,brave,duckduckgo,startpage,wikipedia"`
	DefaultProvider  string   `env:"DEFAULT_PROVIDER, default=bing"`
	UserAgent        string   `env:"USER_AGENT, default=Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36"`
}

var config Config

func loadConfig() {
	ctx := context.Background()
	if err := envconfig.Process(ctx, &config); err != nil {
		log.Fatal(err)
	}
	if config.BaseURL == "" {
		config.BaseURL = fmt.Sprintf("http://127.0.0.1:%d", config.Port)
	}
}
