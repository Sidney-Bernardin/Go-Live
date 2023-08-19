FROM golang:latest AS builder

WORKDIR /app
COPY . .

# Download go-modules and build the app.
RUN go mod download
RUN go build -o service .

# =====

FROM golang:latest

WORKDIR /app
COPY --from=builder /app/service .
