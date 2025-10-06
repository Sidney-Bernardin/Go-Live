FROM golang:1.25.1-alpine AS build

WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o app .

# =====

FROM scratch

COPY --from=build /app/app .

ENTRYPOINT ["./app"]
