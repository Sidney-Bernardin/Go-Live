FROM golang:latest AS builder

# Copy over the source code.
WORKDIR /app
COPY . .

# Download go-modules and build the app.
RUN go mod download
RUN go build -o service .

# ============================================================================

FROM golang:latest

# Copy over the binary.
WORKDIR /app
COPY --from=builder /app/service .
