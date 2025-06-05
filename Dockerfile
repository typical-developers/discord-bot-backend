FROM golang:latest AS builder

ARG POSTGRES_URL

# Set up the build environment
RUN mkdir -p /app
WORKDIR /app

COPY . /app/

# Install Golang CLI tools
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Handle schema generation and db migrations
RUN swag init --parseDependency -dir ./api --output ./internal/docs
RUN sqlc generate --file internal/db/sqlc.yaml
RUN migrate -path internal/db/migrations -database "$POSTGRES_URL" up

# Build the final app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o bin/api-backend

CMD ["./bin/api-backend"]