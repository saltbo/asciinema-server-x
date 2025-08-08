# --- Build frontend ---
FROM node:20-alpine AS web-build
WORKDIR /web
# When building with context at repo root, copy from ./web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web .
RUN npm run build

# --- Build backend ---
FROM golang:1.22-alpine AS builder
WORKDIR /app
# Copy Go module files from server/ when context is repo root
COPY server/go.mod server/go.sum ./
RUN go mod download
# Copy server source
COPY server .
RUN CGO_ENABLED=0 GOOS=linux go build -o /asciinema-server-x ./cmd/server/main.go

# --- Final stage ---
FROM alpine:3.20 AS prepare
RUN mkdir -p /data

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /asciinema-server-x /asciinema-server-x
COPY --from=web-build /web/dist /web-dist
# Pre-create data dir with nonroot ownership
COPY --from=prepare --chown=nonroot:nonroot /data /data
VOLUME ["/data"]
USER nonroot
EXPOSE 8080
ENTRYPOINT ["/asciinema-server-x"]
