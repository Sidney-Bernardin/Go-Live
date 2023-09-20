FROM golang:alpine AS build

WORKDIR /app

# Copy over go.mod and go.sum files to download dependencies.
COPY go.* .
RUN go mod download 

# Copy over the source code to build it.
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o app . 

# ============================================================================

FROM scratch

# Copy over the built application.
COPY --from=build /app/app .
