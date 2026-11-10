# One image: the Go server plus the built Vue app it serves.

# --- frontend -------------------------------------------------------------
FROM oven/bun:1-alpine AS ui
WORKDIR /ui
COPY gamevote-ui/package.json gamevote-ui/bun.lock ./
RUN bun install --frozen-lockfile
COPY gamevote-ui/ ./
# The generated API client is committed, so the spec isn't needed here.
RUN bun --bun run build-only

# --- backend --------------------------------------------------------------
FROM golang:1.25-alpine AS api
WORKDIR /src
COPY gamevote-api-go/go.mod gamevote-api-go/go.sum ./
RUN go mod download
COPY gamevote-api-go/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /gamevote ./cmd/api

# --- runtime --------------------------------------------------------------
FROM alpine:3
# CA certificates are needed for the Steam store search (HTTPS).
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 gamevote
WORKDIR /app
COPY --from=api /gamevote ./gamevote
COPY --from=ui /ui/dist ./web
USER gamevote
ENV PORT=8080 STATIC_DIR=/app/web
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s \
  CMD wget -q -O /dev/null http://127.0.0.1:8080/api/drinks || exit 1
CMD ["./gamevote"]
