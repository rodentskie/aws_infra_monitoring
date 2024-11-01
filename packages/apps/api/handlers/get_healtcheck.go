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

	healthChecks, err := awsquery.FetchHealthChecks(context.Background(), region)
	if err != nil {
		log.Println(err)
	}

	var data []structs.HealthCheckMetric
	for _, hc := range healthChecks {
		metric, err := awsquery.GetRoute53Metrics(context.Background(), region, hc.ID)
		if err != nil {
			log.Println(err)
		}

		if hc.Metrics == nil {
			hc.Metrics = &metric
		}

		data = append(data, structs.HealthCheckMetric{
			HealthCheck: hc,
		})
	}

	res := structs.HealthCheckMetricResponse{
		Data: data,
	}

	j, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
