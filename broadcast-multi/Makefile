.PHONY: build
build:
	@go build -o bin/node cmd/main.go

.PHONY: test
test: build
	@../../maelstrom/maelstrom test -w broadcast --bin bin/node --node-count 5 --time-limit 20 --rate 10