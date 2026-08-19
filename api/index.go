package handler

import (
	"net/http"
	"sync"

	site "mattjs.me/internal/site"
)

var (
	app     http.Handler
	appErr  error
	appOnce sync.Once
)

// Handler is the Vercel Go Function entry point.
func Handler(w http.ResponseWriter, r *http.Request) {
	appOnce.Do(func() {
		app, appErr = site.NewApp()
	})
	if appErr != nil {
		http.Error(w, "failed to initialise application", http.StatusInternalServerError)
		return
	}
	app.ServeHTTP(w, r)
}
