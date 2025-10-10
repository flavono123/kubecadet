# syntax=docker/dockerfile:1

FROM golang:1.22 AS builder
WORKDIR /src

COPY go.mod ./
COPY go.sum ./

# Download dependencies early for better caching. Will succeed once go.sum is populated.
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download || true

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /workspace/app ./cmd/otel-app

FROM gcr.io/distroless/base-debian12:nonroot
USER nonroot:nonroot
WORKDIR /app
COPY --from=builder /workspace/app ./otel-app
EXPOSE 8080
ENTRYPOINT ["./otel-app"]
