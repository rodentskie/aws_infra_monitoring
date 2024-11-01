package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	awsquery "packages/library/go/aws_query"
	"packages/library/go/env"
	"packages/library/go/structs"
)

func GetHealthCheckRequestHandler(w http.ResponseWriter, r *http.Request) {
	defaultRegion := env.GetEnv("REGION", "ap-southeast-2")

	urlQuery := r.URL.Query()
	region := urlQuery.Get("region")
	if region == "" {
		region = defaultRegion
	}

	err := awsquery.GetRoute53Metrics(context.Background(), region, "1c90567f-a5aa-40e2-a8fc-4d7cd0ea5d24")
	if err != nil {
		log.Println(err)
	}

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
