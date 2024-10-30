package handlers

import (
	"encoding/json"
	"net/http"
	"packages/library/structs"
)

func IndexRequestHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes := structs.Response{
		Message: "Welcome to AWS infra monitoring API.",
	}

	j, err := json.Marshal(bodyBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
