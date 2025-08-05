# Build stage
FROM golang:1.22.2 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y yt-dlp libreoffice ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/main .
# ❌ DON'T COPY .env
EXPOSE 8080
CMD ["./main"]
