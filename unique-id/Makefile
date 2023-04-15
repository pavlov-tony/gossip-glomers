.PHONY: build
build:
	@go build -o bin/node cmd/main.go

.PHONY: test
test: build
	@../../maelstrom/maelstrom test -w unique-ids --bin bin/node --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
