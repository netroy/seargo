ARG GO_VERSION=1.24.3

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src
RUN apk add --no-cache make
COPY go.mod go.sum Makefile /src
RUN make setup
COPY . /src
RUN make build

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /src/personal-search /app/personal-search
COPY templates /app/templates
EXPOSE 8080
ENTRYPOINT ["/app/personal-search"]
