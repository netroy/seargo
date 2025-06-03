setup:
	go mod download

dev:
	go run .

build:
	go build -ldflags "-s -w" -o "personal-search"

build-docker:
	docker build -t personal-search:dev .

test:
	go test -v .

cover:
	go test -coverprofile cover.prof . && covreport -i cover.prof -o cover.html -cutlines 70,40 && xdg-open cover.html

format:
	gofmt -w .
