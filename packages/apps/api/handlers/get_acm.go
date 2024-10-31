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

func GetCertificatesRequestHandler(w http.ResponseWriter, r *http.Request) {
	defaultRegion := env.GetEnv("REGION", "ap-southeast-2")

	urlQuery := r.URL.Query()
	region := urlQuery.Get("region")
	if region == "" {
		region = defaultRegion
	}

	cert, err := awsquery.GetACMCertificates(context.Background(), region)
	if err != nil {
		log.Println(err)
	}

	res := structs.ACMCertificateResponse{
		Data: cert,
	}

	j, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
