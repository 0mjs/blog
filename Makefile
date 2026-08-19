.PHONY: install generate css build test run dev dev/build

install:
	go mod download
	npm install

generate:
	go tool templ generate

css:
	npm run css

build: generate css
	mkdir -p bin
	go build -o bin/blog .

test: generate
	go test ./...

run: generate css
	go run .

dev:
	air

dev/build:
	mkdir -p tmp
	go tool templ generate
	npm run css
	go build -o tmp/blog .
