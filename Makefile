.PHONY: install generate css build test run

install:
	go mod download
	npm install

generate:
	go tool templ generate

css:
	npm run css

build: generate
	mkdir -p bin
	go build -o bin/blog .

test: generate
	go test ./...

run: generate
	go run .
