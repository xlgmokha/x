# minit

A minimal Procfile-based process supervisor for Docker containers.

## Features

- **Zero dependencies**: Works in scratch containers without shells or external tools
- **Procfile support**: Run multiple processes from a single configuration file
- **Signal forwarding**: Graceful shutdown with SIGTERM propagation
- **Process groups**: Uses `setpgid` for proper signal handling
- **Clean logs**: Direct stdout/stderr streaming without buffering

## Usage

```bash
# Run (reads Procfile from current directory)
./minit
```

## Procfile Format

```
web: ./server -port 8080
worker: ./worker -queue tasks
redis: redis-server --port 6379
scheduler: ./scheduler -interval 60s
```

## Installation

```dockerfile
FROM golang:alpine AS minit-builder
RUN go install github.com/xlgmokha/minit@latest
```

## Docker Example
Combine minit with [dumb-init](https://github.com/Yelp/dumb-init) for proper PID 1 signal handling:

```dockerfile
# syntax=docker/dockerfile:1

# Build stage for minit
FROM golang:alpine AS minit-builder
RUN go install github.com/xlgmokha/minit@latest

# Build stage for dumb-init
FROM debian:bookworm-slim AS dumb-init-builder
RUN apt-get update && apt-get install -y wget && \
    wget -O /usr/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.5/dumb-init_1.2.5_x86_64 && \
    chmod +x /usr/bin/dumb-init

# Final stage
FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=minit-builder /go/bin/minit /bin/minit
COPY --from=dumb-init-builder /usr/bin/dumb-init /usr/bin/dumb-init
COPY Procfile /Procfile
COPY your-app /your-app

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/bin/minit"]
```

### Why use dumb-init?
- **PID 1 responsibilities**: Zombie reaping, proper signal handling
- **Graceful shutdown**: Ensures all processes terminate cleanly
- **Container compatibility**: Works with all container runtimes
- **Security**: Runs as non-root user with distroless base

## License

MIT
