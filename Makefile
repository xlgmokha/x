.PHONY: clean setup build test

syslogs:
	@go build -o syslogs ./cmd/syslogs/main.go

which:
	@go build -o which ./cmd/which/main.go

clean:
	@rm -f syslogs which

setup:
	@mise install
	@go install tool

build: syslogs which

test:
	@go clean -testcache
	@go test -shuffle=on ./...
