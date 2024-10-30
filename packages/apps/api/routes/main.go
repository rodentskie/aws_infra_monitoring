package routes

import (
	"net/http"
	"packages/apps/api/handlers"
)

func MainRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handlers.IndexRequestHandler)
}
