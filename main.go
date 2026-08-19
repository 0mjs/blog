package main

import (
	"log"
	"os"
)

func main() {
	app, err := newApp()
	if err != nil {
		log.Fatal(err)
	}
	address := os.Getenv("ADDRESS")
	if address == "" {
		address = ":3000"
	}
	log.Printf("Zinc blog listening on http://localhost%s", address)
	log.Fatal(app.Listen(address))
}
