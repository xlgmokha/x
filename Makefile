.PHONY: build clean test vet syslogs which minit http-server proxy-server

syslogs:
	@mkdir -p bin
	@go build -o bin/syslogs ./cmd/syslogs

which:
	@mkdir -p bin
	@go build -o bin/which ./cmd/which

minit:
	@mkdir -p bin
	@go build -o bin/minit ./cmd/minit

http-server:
	@mkdir -p bin
	@go build -o bin/http-server ./cmd/http-server

proxy-server:
	@mkdir -p bin
	@go build -o bin/proxy-server ./cmd/proxy-server

build: syslogs which minit http-server proxy-server

clean:
	@rm -rf bin

vet:
	@go vet ./...

test:
	@go clean -testcache
	@go test -race -shuffle=on ./...
