package datafetcher

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// HeartbeatsFetcher is the implementation of DataFetcher that uses a queryExecutor to perform a query
//and obtain the Heartbeats.
type HeartbeatsFetcher struct {
	queryExecutor cloudwatch.QueryExecutor
}

// NewHeartbeatsFetcher creates a new HeartbeatsFetcher implementationm
func NewHeartbeatsFetcher(queryExecutor cloudwatch.QueryExecutor) HeartbeatsFetcher {
	return HeartbeatsFetcher{queryExecutor}
}

// GetLogGroupName returns the appropiate Log group in AWS CloudWatch
func (heartbeatsFetcher HeartbeatsFetcher) GetLogGroupName() (logGroupName string) {
	return "/aws/lambda/AWSUpload"
}

// createQueryTemplate returns a new query template depending on the productNumber and serialNumber parameters.
// The resulting query template will be used by the queryExecutor to obtain the Heartbeats.
func (heartbeatsFetcher HeartbeatsFetcher) CreateQueryTemplate(productNumber, serialNumber string) (queryTemplateString string) {
	if productNumber != "" && serialNumber != "" {
		return `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date) and fields.topic = "heartbeat" and fields.ProductNumber="{{.productNumber}}" and fields.SerialNumber="{{.serialNumber}}"
								| sort @timestamp asc
								| limit 10000`
	}
	if productNumber != "" {
		return `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date) and fields.topic = "heartbeat" and fields.ProductNumber="{{.productNumber}}"
								| sort @timestamp asc
								| limit 10000`
	}

	return `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date) and fields.topic = "heartbeat"
								| sort @timestamp asc
								| limit 10000`
}

// FetchData obtains the uploaded Heartbeats depending on requestQueryParams.
// The method basically creates a variable insightsQueryParams and then calls a queryExecutor
// to perform the query. It returns the result and an error, if any.
func (heartbeatsFetcher HeartbeatsFetcher) FetchData(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	insightsQueryParams, err := createInsightsQueryParams(requestQueryParams, heartbeatsFetcher)
	if err != nil {
		return nil, err
	}

	result, err := heartbeatsFetcher.queryExecutor.ExecuteQuery(insightsQueryParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}
