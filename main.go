package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/aaronlee232/redis-gui-tester/api"
	ui "github.com/aaronlee232/redis-gui-tester/ui/templates"
)

func main() {
	fmt.Println("Running redis-gui-tester")

	// Start Server
	addr := ":3000"
	apiRouter := api.NewAPIRouter()
	frontendRouter := NewFrontendRouter()

	mainRouter := http.NewServeMux()
	mainRouter.Handle("/api/", http.StripPrefix("/api", apiRouter))
	mainRouter.Handle("/", frontendRouter)

	log.Println("Starting server on port", addr)
	if err := http.ListenAndServe(addr, mainRouter); err != nil {
		log.Fatal(err)
	}
}

func NewFrontendRouter() *http.ServeMux {
	// Serve frontend pages
	component := ui.Layout("Redis GUI Tester")

	r := http.NewServeMux()

	// Serve static files (CSS, JS) under /static/
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))))
	r.Handle("/", templ.Handler(component))

	return r
	// fmt.Println("Listening on :3000")
	// log.Fatal(http.ListenAndServe(":3000", nil))
}
