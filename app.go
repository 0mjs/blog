package main

import (
	"github.com/0mjs/zinc"
	site "mattjs.me/internal/site"
)

func newApp() (*zinc.App, error) {
	return site.NewApp()
}
