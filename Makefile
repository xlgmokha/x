.PHONY: build clean test vet syslogs which

syslogs:
	@mkdir -p bin
	@go build -o bin/syslogs ./cmd/syslogs

which:
	@mkdir -p bin
	@go build -o bin/which ./cmd/which

build: syslogs which

clean:
	@rm -rf bin

vet:
	@go vet ./...

test:
	@go clean -testcache
	@go test -race -shuffle=on ./...
