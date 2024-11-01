package structs

import "github.com/aws/aws-sdk-go-v2/service/route53/types"

type HealthCheckMetrics struct {
	Time   string  `json:"time"`
	Status float64 `json:"status"`
}

type HealthCheck struct {
	ID                string                  `json:"id"`
	Type              types.HealthCheckType   `json:"type"`
	HealthCheckConfig types.HealthCheckConfig `json:"healthCheckConfig"`
	Metrics           *[]HealthCheckMetrics   `json:"metrics,omitempty"`
}

type HealthCheckMetric struct {
	HealthCheck HealthCheck `json:"healthCheck"`
}

type HealthCheckMetricResponse struct {
	Data []HealthCheckMetric `json:"data"`
}
