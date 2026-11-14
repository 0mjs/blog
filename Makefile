.PHONY: install generate css build test run dev dev/build write

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

# The local writing app: edit posts with a live preview, add images, publish. http://127.0.0.1:4000
write: generate
	go run ./cmd/editor

dev/build:
	mkdir -p tmp
	go tool templ generate
	npm run css
	go build -o tmp/blog .
