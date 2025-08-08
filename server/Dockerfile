# --- Build frontend ---
FROM node:20-alpine AS web-build
WORKDIR /web
COPY ../web/package.json ../web/package-lock.json ./
RUN npm ci --omit=dev
COPY ../web .
RUN npm run build

# --- Build backend ---
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /asciinema-server-x ./cmd/server/main.go

# --- Final stage ---
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /asciinema-server-x /asciinema-server-x
COPY --from=web-build /web/dist /web-dist
# Create data dir (if not using PVC)
RUN mkdir -p /data && chown nonroot:nonroot /data
VOLUME ["/data"]
USER nonroot
EXPOSE 8080
ENTRYPOINT ["/asciinema-server-x"]
