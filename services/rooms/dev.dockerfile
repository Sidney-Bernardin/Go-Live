FROM golang:alpine AS base

WORKDIR /app

# Copy over go.mod and go.sum files to download dependencies.
COPY go.* .
RUN go mod download 

# Copy over the rest of the source code.
COPY . .

# ============================================================================

FROM base AS build

# Build the application.
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o app . 

# ============================================================================

FROM base AS test

# Test the application.
RUN --mount=type=cache,target=/root/.cache/go-build \
    go test -v .

# ============================================================================

FROM scratch

# Copy over the built application.
COPY --from=build /app/app .
