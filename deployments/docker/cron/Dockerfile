FROM golang:latest AS builder

# Set up the build environment
RUN mkdir -p /app
WORKDIR /app

COPY . /app/

# Build the final app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o ./bin/cron ./services/cron/cmd/main.go

CMD ["./bin/cron"]