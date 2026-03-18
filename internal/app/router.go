package app

import (
	"net/http"

	"github.com/aaronlee232/redis-gui-tester/internal/database"
	"github.com/aaronlee232/redis-gui-tester/internal/middleware"
	"github.com/aaronlee232/redis-gui-tester/internal/scenario"
	"github.com/aaronlee232/redis-gui-tester/internal/tester"
)

func NewRouter(repo *database.Registry) *http.ServeMux {
	mux := http.NewServeMux()

	scenarioHandler := scenario.NewHandler(repo)
	testerHandler := tester.NewHandler(repo)

	mux.Handle("/api/scenario/", http.StripPrefix("/api/scenario", middleware.StripTrailingSlash(scenarioHandler.RegisterRoutes())))
	mux.Handle("/api/tester/", http.StripPrefix("/api/tester", middleware.StripTrailingSlash(testerHandler.RegisterRoutes())))

	return mux
}
