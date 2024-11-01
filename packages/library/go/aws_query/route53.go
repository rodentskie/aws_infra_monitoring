package awsquery

import (
	"context"
	"fmt"
	"packages/library/go/structs"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	ct "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func GetRoute53Metrics(ctx context.Context, region, healthCheckID string) ([]structs.HealthCheckMetrics, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	cwClient := cloudwatch.NewFromConfig(cfg)
	r53Client := route53.NewFromConfig(cfg)

	// Get current time and time 1 week ago
	endTime := time.Now()
	startTime := endTime.Add(-7 * 24 * time.Hour)

	_, err = r53Client.GetHealthCheck(ctx, &route53.GetHealthCheckInput{
		HealthCheckId: &healthCheckID,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting health check: %v", err)
	}

	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Route53"),
		MetricName: aws.String("HealthCheckStatus"),
		Dimensions: []ct.Dimension{
			{
				Name:  aws.String("HealthCheckId"),
				Value: aws.String(healthCheckID),
			},
		},
		StartTime:  aws.Time(startTime),
		EndTime:    aws.Time(endTime),
		Period:     aws.Int32(7200), // 2 hour periods
		Statistics: []ct.Statistic{ct.StatisticAverage},
	}

	result, err := cwClient.GetMetricStatistics(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting metrics: %v", err)
	}

	if len(result.Datapoints) == 0 {
		return nil, fmt.Errorf("error no data: %v", err)
	}

	// Sort datapoints by timestamp
	sort.Slice(result.Datapoints, func(i, j int) bool {
		return result.Datapoints[i].Timestamp.Before(*result.Datapoints[j].Timestamp)
	})

	var checks []structs.HealthCheckMetrics
	for _, dp := range result.Datapoints {
		checks = append(checks, structs.HealthCheckMetrics{
			Time:   dp.Timestamp.Format("2006-01-02 15:04:05"),
			Status: *dp.Average,
		})
	}
	return checks, nil
}

func FetchHealthChecks(ctx context.Context, region string) ([]structs.HealthCheck, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %v", err)
	}

	client := route53.NewFromConfig(cfg)

	input := &route53.ListHealthChecksInput{}
	paginator := route53.NewListHealthChecksPaginator(client, input)

	var healthChecks []structs.HealthCheck

	// Paginate through health checks
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list health checks: %v", err)
		}

		// Process each health check
		for _, hc := range output.HealthChecks {
			if hc.HealthCheckConfig != nil {
				healthCheck := structs.HealthCheck{
					ID:                *hc.Id,
					Type:              hc.HealthCheckConfig.Type,
					HealthCheckConfig: *hc.HealthCheckConfig,
				}
				healthChecks = append(healthChecks, healthCheck)
			}
		}
	}

	return healthChecks, nil
}
