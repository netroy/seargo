# Seargo

This is a minimal metasearch engine for my personal usage.

While many mature metasearch engines exist, I prefer services that I self-host to be as lean as possible, which most available metasearch engines aren't.

The docker image for this project, and the memory usage under casual load is under **20MBs**, and that is important to me.

## Usage

- run `docker run -d --rm -p 8080:8080 ghcr.io/netroy/seargo` to start the container.
- visit `http://127.0.0.1:8080/search?q=[[YOUR_QUERY]]` to search for something.


## Providers

Currently only duckduckgo is being used, but I plan to add more providers soon.
