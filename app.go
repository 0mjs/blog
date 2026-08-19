package main

import (
	"github.com/0mjs/zinc"
	"mattjs.me/site"
)

func newApp() (*zinc.App, error) {
	return site.NewApp()
}
