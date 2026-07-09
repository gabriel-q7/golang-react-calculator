# syntax=docker/dockerfile:1

# ---- Stage 1: build the frontend ------------------------------------------
FROM node:22-alpine AS frontend-builder
WORKDIR /src/frontend
COPY apps/frontend/package.json apps/frontend/package-lock.json* ./
RUN npm ci
COPY apps/frontend/ ./
RUN npm run build

# ---- Stage 2: build the Go binary, with the frontend embedded -------------
FROM golang:1.23-alpine AS backend-builder
WORKDIR /src/backend
COPY apps/backend/go.mod ./
RUN go mod download
COPY apps/backend/ ./
# Embedded via go:embed in internal/web/embed.go — must exist before `go build`.
COPY --from=frontend-builder /src/frontend/dist/ ./internal/web/dist/
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# ---- Stage 3: minimal runtime ----------------------------------------------
FROM alpine:3.20 AS runtime
RUN apk add --no-cache ca-certificates && \
    adduser -D -u 10001 app
COPY --from=backend-builder /out/server /usr/local/bin/server
USER app
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["server"]
