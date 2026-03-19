# Use Go 1.25 (or your preferred version)
FROM golang:1.25-alpine

WORKDIR /app

# Copy Go modules (optional if you have go.mod)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY getMetrics.go .

# Build the metrics server
RUN go build -o getMetrics getMetrics.go

EXPOSE 9100

CMD ["./getMetrics"]
