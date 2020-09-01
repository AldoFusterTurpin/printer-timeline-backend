package datafetcher

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// OpenXmlsFetcher is the implementation of DataFetcher that uses a queryExecutor to perform a query
//and obtain the OpenXMls.
type OpenXmlsFetcher struct {
	queryExecutor cloudwatch.QueryExecutor
}

// NewOpenXmlsFetcher creates a new OpenXmlsFetcher implementationm
func NewOpenXmlsFetcher(queryExecutor cloudwatch.QueryExecutor) OpenXmlsFetcher {
	return OpenXmlsFetcher{queryExecutor}
}

// GetLogGroupName returns the appropiate Log group in AWS CloudWatch
func (openXmlsFetcher OpenXmlsFetcher) GetLogGroupName() (logGroupName string) {
	return "/aws/lambda/AWSUpload"
}

// createQueryTemplate returns a new query template depending on the productNumber and serialNumber parameters.
// The resulting query template will be used by the queryExecutor to obtain the OpenXMls.
func (openXmlsFetcher OpenXmlsFetcher) CreateQueryTemplate(productNumber, serialNumber string) (queryTemplateString string) {
	if productNumber != "" && serialNumber != "" {
		return `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date) and fields.topic != "heartbeat" and fields.ProductNumber="{{.productNumber}}" and fields.SerialNumber="{{.serialNumber}}"
								| sort @timestamp asc
								| limit 10000`
	}
	if productNumber != "" {
		return `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date) and fields.topic != "heartbeat" and fields.ProductNumber="{{.productNumber}}"
								| sort @timestamp asc
								| limit 10000`
	}

	return `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date) and fields.topic != "heartbeat"
								| sort @timestamp asc
								| limit 10000`
}

// FetchData obtains the uploaded OpenXml depending on requestQueryParams.
// The method basically creates a variable insightsQueryParams and then calls a queryExecutor
// to perform the query. It returns the result and an error, if any.
func (openXmlsFetcher OpenXmlsFetcher) FetchData(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	insightsQueryParams, err := createInsightsQueryParams(requestQueryParams, openXmlsFetcher)
	if err != nil {
		return nil, err
	}

	result, err := openXmlsFetcher.queryExecutor.ExecuteQuery(insightsQueryParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}
