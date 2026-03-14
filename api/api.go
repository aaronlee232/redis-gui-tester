package api

import (
	"net/http"
	"strings"

	"github.com/aaronlee232/redis-gui-tester/api/services/scenario"
	"github.com/aaronlee232/redis-gui-tester/api/services/tester"
)

// stripTrailingSlash removes trailing slashes and normalizes empty paths to "/"
func stripTrailingSlash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "/" && strings.HasSuffix(path, "/") {
			path = strings.TrimSuffix(path, "/")
		}
		if path == "" {
			path = "/"
		}
		r.URL.Path = path
		h.ServeHTTP(w, r)
	})
}

func NewAPIRouter() *http.ServeMux {
	// Register Routes
	scenarioHandler := scenario.NewHandler()
	scenarioRouter := stripTrailingSlash(scenarioHandler.RegisterRoutes())

	testerHandler := tester.NewHandler()
	testerRouter := stripTrailingSlash(testerHandler.RegisterRoutes())

	// Attach to base URLs
	r := http.NewServeMux()
	r.Handle("/scenario/", http.StripPrefix("/scenario", scenarioRouter))
	r.Handle("/tester/", http.StripPrefix("/tester", testerRouter))

	return r
}
