# ---- Stage 1: Build Go backend ----
FROM golang:1.21-alpine AS go-builder
WORKDIR /build
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# ---- Stage 2: Build Next.js frontend ----
FROM node:20-alpine AS node-builder
WORKDIR /build
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci
COPY frontend/ .
RUN mkdir -p public && npm run build

# ---- Stage 3: Runtime ----
FROM node:20-alpine
RUN apk add --no-cache ca-certificates tzdata nginx

WORKDIR /app

# Go binary (statically compiled, runs on any alpine)
COPY --from=go-builder /build/server ./server

# Next.js standalone bundle
COPY --from=node-builder /build/.next/standalone ./frontend
COPY --from=node-builder /build/.next/static     ./frontend/.next/static
COPY --from=node-builder /build/public           ./frontend/public

# nginx template + startup script
COPY nginx.conf.template /etc/nginx/nginx.conf.template
COPY start.sh            ./start.sh
RUN chmod +x ./start.sh

EXPOSE 8080
CMD ["./start.sh"]
