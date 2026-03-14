package ui

import (
	"net/http"

	"github.com/a-h/templ"
	ui "github.com/aaronlee232/redis-gui-tester/ui/templates"
)

func NewFrontendRouter() *http.ServeMux {
	// Serve frontend pages
	rootLayout := ui.Layout("Redis GUI Tester")

	r := http.NewServeMux()

	// Serve static files (CSS, JS) under /static/
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))))
	r.Handle("/", templ.Handler(rootLayout))

	return r
}
