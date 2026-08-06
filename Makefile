.PHONY: build ci clean fmt fmt-check install test tidy vet \
        syslogs which minit http-server proxy-server

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

install:
	@go install ./cmd/...

clean:
	@rm -rf bin

fmt:
	@gofmt -s -l -w .

fmt-check:
	@unformatted=$$(gofmt -s -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt needed:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

tidy:
	@go mod tidy

vet:
	@go vet ./...

test:
	@go clean -testcache
	@go test -race -shuffle=on ./...

ci: fmt-check vet test
