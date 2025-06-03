# Seargo

This is a minimal metasearch engine for my personal usage.

## Usage

- run `docker run -d --rm -p 8080:8080 ghcr.io/netroy/seargo` to start the container.
- visit `http://127.0.0.1:8080/search?q=[[YOUR_QUERY]]` to search for something.

## Providers

Currently only Duckduckgo, Brave, and English wikipedia are being used, but I plan to add more providers soon.
