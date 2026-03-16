# Stage 1: build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy module file first; source files are needed for go mod tidy
COPY go.mod ./
COPY *.go ./

# go mod tidy downloads dependencies and generates go.sum in one step
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /search-service .

# Stage 2: minimal runtime image
FROM alpine:3.19

# ca-certificates is required for TLS calls to the Gemini API
RUN apk add --no-cache ca-certificates

COPY --from=builder /search-service /search-service

EXPOSE 8083

ENTRYPOINT ["/search-service"]
