# Seargo

This is a minimal metasearch engine for my personal usage. It supports searching via browsers, as well as via AI agents over MCP.

## Usage

- run `docker run -d --rm -p 8080:8080 ghcr.io/netroy/seargo` to start the container.
- visit `http://127.0.0.1:8080/search?q=[[YOUR_QUERY]]` to search for something.
- Use the following URLs for MCP
  - SSE: `http://127.0.0.1:8080/sse`
  - Streamable HTTP: `http://127.0.0.1:8080/mcp`

## Providers

Currently only Bing, Brave, DuckDuckGo, Startpage, and English wikipedia are being used, but I plan to add more providers as I need them.
