build:
	@go build  -o bin/blockchain

run: build
	@./bin/docker

test:
	@go test -v ./...