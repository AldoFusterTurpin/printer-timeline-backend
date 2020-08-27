package datafetcher

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudwatch"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/queryparams"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// HeartbeatsFetcher obtains the Uploaded Heartbeats based on request query parameters.
type HeartbeatsFetcher interface {
	FetchData(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error)
	GetUploadedHeartbeats(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error)
}

// HeartbeatsFetcherImpl is the implementation of HeartbeatsFetcher that uses a queryExecutor to perform a query
//and obtain the Heartbeats.
type HeartbeatsFetcherImpl struct {
	queryExecutor cloudwatch.QueryExecutor
}

// NewHeartbeatsFetcherImpl creates a new HeartbeatsFetcher implementationm
func NewHeartbeatsFetcherImpl(queryExecutor cloudwatch.QueryExecutor) HeartbeatsFetcher {
	return HeartbeatsFetcherImpl{queryExecutor}
}

// createQueryTemplate returns a new query template depending on the productNumber and serialNumber parameters.
// The resulting query template will be used by the queryExecutor to obtain the Heartbeats.
func (heartbeatsFetcherImpl HeartbeatsFetcherImpl) createQueryTemplate(productNumber, serialNumber string) (queryTemplateString string) {
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

// createInsightsQueryParams creates InsightQueryParameters based on requestQueryParams.
// The returned InsightQueryParameters will be used by a QueryExecutor to execute the query. It also returns an error, if any.
func (heartbeatsFetcherImpl HeartbeatsFetcherImpl) createInsightsQueryParams(requestQueryParams map[string]string) (insightsQueryParams cloudwatch.InsightsQueryParams, err error) {
	startTime, endTime, err := queryparams.ExtractTimeRange(requestQueryParams)
	if err != nil {
		return
	}

	productNumber, serialNumber, err := queryparams.ExtractPrinterInfo(requestQueryParams)
	if err != nil {
		return
	}

	queryTemplate := heartbeatsFetcherImpl.createQueryTemplate(productNumber, serialNumber)

	mapValues := map[string]string{
		"productNumber": productNumber,
		"serialNumber":  serialNumber,
	}

	queryToExecute, err := cloudwatch.CreateQuery(queryTemplate, mapValues)
	if err != nil {
		return
	}

	insightsQueryParams = cloudwatch.InsightsQueryParams{
		StartTimeEpoch: startTime.Unix(),
		EndTimeEpoch:   endTime.Unix(),
		LogGroupName:   "/aws/lambda/AWSUpload",
		Query:          queryToExecute,
	}
	return insightsQueryParams, nil
}

// GetUploadedHeartbeats obtains the uploaded Heartbeats depending on requestQueryParams.
// The method basically creates a variable insightsQueryParams and then calls a queryExecutor
// to perform the query. It returns the result and an error, if any.
func (heartbeatsFetcherImpl HeartbeatsFetcherImpl) GetUploadedHeartbeats(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	insightsQueryParams, err := heartbeatsFetcherImpl.createInsightsQueryParams(requestQueryParams)
	if err != nil {
		return nil, err
	}

	result, err := heartbeatsFetcherImpl.queryExecutor.ExecuteQuery(insightsQueryParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (heartbeatsFetcherImpl HeartbeatsFetcherImpl) FetchData(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	return heartbeatsFetcherImpl.GetUploadedHeartbeats(requestQueryParams)
}
