package main

import (
	"log"

	"blog.0mjs.dev/site"
)

func main() {
	app := site.NewApp()

	log.Fatal(app.Listen(":3000"))
}
