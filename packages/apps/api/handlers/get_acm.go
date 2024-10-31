package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	awsquery "packages/library/go/aws_query"
	"packages/library/go/structs"
)

func GetCertificatesRequestHandler(w http.ResponseWriter, r *http.Request) {
	cert, err := awsquery.GetACMCertificates(context.Background(), "ap-southeast-2")
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
