package app

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/aaronlee232/redis-gui-tester/internal/database"
	"github.com/aaronlee232/redis-gui-tester/internal/middleware"
	"github.com/aaronlee232/redis-gui-tester/internal/scenario"
	"github.com/aaronlee232/redis-gui-tester/internal/tester"
	"github.com/aaronlee232/redis-gui-tester/internal/ui"
)

func NewRouter(repo *database.Registry) *http.ServeMux {
	mux := http.NewServeMux()

	// 1. Static Assets
	fileServer := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// 2. API Routes (mounted with middleware)
	scenarioHandler := scenario.NewHandler(repo)
	testerHandler := tester.NewHandler(repo)

	// Attach to base URLs
	mux.Handle("/api/scenario/", http.StripPrefix("/api/scenario", middleware.StripTrailingSlash(scenarioHandler.RegisterRoutes())))
	mux.Handle("/api/tester/", http.StripPrefix("/api/tester", middleware.StripTrailingSlash(testerHandler.RegisterRoutes())))

	// Serve static files (CSS, JS) under /static/
	rootLayout := ui.Layout("Redis GUI Tester")
	mux.Handle("/", templ.Handler(rootLayout))

	return mux
}
