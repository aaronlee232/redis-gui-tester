package scenario

import (
	"fmt"
	"net/http"
)

func GetScenario(w http.ResponseWriter, r *http.Request) {

}

func GetAllScenarios(w http.ResponseWriter, r *http.Request) {
	// 1. Set a custom header
	w.Header().Set("Content-Type", "text/plain")

	// 2. Set the HTTP status code (must be called before w.Write)
	w.WriteHeader(http.StatusOK)

	// 3. Write the response body
	fmt.Fprintf(w, "Hello, World!")
}

func CreateScenario(w http.ResponseWriter, r *http.Request) {

}

func UpdateScenario(w http.ResponseWriter, r *http.Request) {

}

func DeleteScenario(w http.ResponseWriter, r *http.Request) {

}
