package awsquery

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func GetRoute53Metrics(ctx context.Context, region, healthCheckID string) error {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return err
	}

	cwClient := cloudwatch.NewFromConfig(cfg)
	r53Client := route53.NewFromConfig(cfg)

	// Get current time and time 1 week ago
	endTime := time.Now()
	startTime := endTime.Add(-7 * 24 * time.Hour)

	_, err = r53Client.GetHealthCheck(context.TODO(), &route53.GetHealthCheckInput{
		HealthCheckId: &healthCheckID,
	})
	if err != nil {
		return fmt.Errorf("error getting health check: %v", err)
	}

	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Route53"),
		MetricName: aws.String("HealthCheckStatus"),
		Dimensions: []types.Dimension{
			{
				Name:  aws.String("HealthCheckId"),
				Value: aws.String(healthCheckID),
			},
		},
		StartTime:  aws.Time(startTime),
		EndTime:    aws.Time(endTime),
		Period:     aws.Int32(3600), // 1 hour periods
		Statistics: []types.Statistic{types.StatisticAverage},
	}

	result, err := cwClient.GetMetricStatistics(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("error getting metrics: %v", err)
	}

	// Print metrics
	fmt.Printf("Metrics for Health Check %s from %s to %s:\n",
		healthCheckID, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	if len(result.Datapoints) == 0 {
		fmt.Println("No data points found")
		return nil
	}

	// Sort datapoints by timestamp
	sort.Slice(result.Datapoints, func(i, j int) bool {
		return result.Datapoints[i].Timestamp.Before(*result.Datapoints[j].Timestamp)
	})

	// Print each datapoint
	for _, dp := range result.Datapoints {
		fmt.Printf("Time: %s, Status: %.2f\n",
			dp.Timestamp.Format("2006-01-02 15:04:05"),
			*dp.Average)
	}

	return nil
}
