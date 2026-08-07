# See docs/docker-notes.md for the full explanation of what's happening
# here and why. Short version: this is a multi-stage build — a fat "builder"
# stage compiles a static Go binary, then a tiny final stage ships only that
# binary, not the Go toolchain or source tree.

# ---- Stage 1: build ----
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Copy just the module files first and download deps before copying the
# rest of the source. Docker caches each instruction as a layer; as long as
# go.mod/go.sum don't change, this layer is reused on rebuilds instead of
# re-downloading dependencies every time you change game code.
COPY go.mod ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 forces a fully static binary (no dynamic C library links),
# which is what lets the final stage be based on a minimal image with
# nothing else installed.
RUN CGO_ENABLED=0 go build -o /out/textfighter ./cmd/textfighter

# ---- Stage 2: run ----
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /out/textfighter .

# Interactive game read from stdin, so this is meant to be run with
# `docker run -it`, not as a detached background service.
ENTRYPOINT ["./textfighter"]
