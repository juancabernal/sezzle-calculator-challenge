# ---- Build stage ----
# golang:1.23-alpine includes the full Go toolchain on a small Alpine
# Linux base — we need this to compile, but we won't ship it.
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Copy only the dependency manifests first. Docker caches each layer;
# as long as go.mod/go.sum don't change, this layer (and the slow
# `go mod download` below) gets reused on future builds, even if the
# application code changes constantly.
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the source and compile.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# ---- Final stage ----
# alpine is a minimal Linux distribution (~5MB) — just enough to run
# a static binary, none of the Go toolchain or source code included.
FROM alpine:latest

WORKDIR /app
COPY --from=builder /server .
COPY migrations ./migrations

EXPOSE 8080

CMD ["./server"]