.PHONY: build
build:
	@go build -o bin/node cmd/main.go

.PHONY: test
test: build
	@../../maelstrom/maelstrom test -w echo --bin bin/node --node-count 1 --time-limit 10
