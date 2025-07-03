# minit

A minimal Procfile-based process supervisor for scratch Docker containers.

## Features

- **No dependencies**: Works in scratch containers without shells
- **Proper signal handling**: Graceful shutdown like dumb-init
- **Clean logs**: Streams stdout/stderr without modifications
- **Process groups**: Ensures all child processes are terminated

## Usage

```bash
# Build
go build -o minit cmd/minit/main.go

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

## Docker Example

```dockerfile
FROM scratch
COPY minit /minit
COPY Procfile /Procfile
COPY your-app /your-app
ENTRYPOINT ["/minit"]
```

## License

MIT
