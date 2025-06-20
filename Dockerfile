ARG GO_VERSION=1.24.4

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src
RUN apk add --no-cache make
COPY go.mod go.sum Makefile /src
RUN make setup
COPY . /src
RUN make build

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /src/seargo /app/seargo
COPY templates /app/templates
COPY static /app/static
EXPOSE 8080
ENTRYPOINT ["/app/seargo"]
