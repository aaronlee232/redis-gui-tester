package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aaronlee232/redis-gui-tester/api"
	"github.com/aaronlee232/redis-gui-tester/ui"
)

func main() {
	fmt.Println("Running redis-gui-tester")

	// Start Server
	addr := ":3000"
	apiRouter := api.NewAPIRouter()
	frontendRouter := ui.NewFrontendRouter()

	mainRouter := http.NewServeMux()
	mainRouter.Handle("/api/", http.StripPrefix("/api", apiRouter))
	mainRouter.Handle("/", frontendRouter)

	log.Println("Starting server on port", addr)
	if err := http.ListenAndServe(addr, mainRouter); err != nil {
		log.Fatal(err)
	}
}
