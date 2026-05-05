# Build stage
FROM golang:1.26 AS builder

WORKDIR /app

# Cache deps first
COPY go.mod go.sum* ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o fleet-app ./main.go


# Runtime stage
FROM debian:bookworm-slim

WORKDIR /app

# Copy built binary
COPY --from=builder /app/fleet-app /app/fleet-app

# Copy runtime files
COPY devices.csv /app/devices.csv

EXPOSE 6733

CMD ["./fleet-app"]