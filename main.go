package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	loadConfig()

	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: mux,
	}

	setupMCPServer(httpServer, mux)
	setupWebServer(mux)

	log.Printf("Server listening on %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
