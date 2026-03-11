package handler

import (
	"net/http"

	"blog.0mjs.dev/site"
)

var app http.Handler = site.NewApp()

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
