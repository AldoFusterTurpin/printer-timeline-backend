package datafetcher

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// RtaFetcher is the implementation of DataFetcher that uses a queryExecutor to perform a query
//and obtain the Rtas.
type RtaFetcher struct {
	queryExecutor cloudwatch.QueryExecutor
}

// NewRtaFetcher creates a new RtaFetcher implementation.
func NewRtaFetcher(queryExecutor cloudwatch.QueryExecutor) RtaFetcher {
	return RtaFetcher{queryExecutor}
}

// GetLogGroupName returns the appropiate Log group in AWS CloudWatch
func (rtasFetcher RtaFetcher) GetLogGroupName() (logGroupName string) {
	return "/aws/lambda/AWSUploadRTA"
}

// CreateQueryTemplate returns a new query template depending on the productNumber and serialNumber parameters.
// The resulting query template will be used by the queryExecutor to obtain the Rtas.
func (rtasFetcher RtaFetcher) CreateQueryTemplate(productNumber, serialNumber string) (queryTemplateString string) {
	if productNumber != "" && serialNumber != "" {
		s1 := "fields @timestamp, `fields.metadata.device-product-number`, `fields.metadata.device-serial-number`, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, `fields.metadata.xml-generator-object-path`"
		s2 := "| filter ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(`fields.metadata.xml-generator-object-path`) and ispresent(`fields.metadata.device-product-number`) and ispresent(`fields.metadata.device-serial-number`) and `fields.metadata.device-product-number`='{{.productNumber}}' and `fields.metadata.device-serial-number`='{{.serialNumber}}'"
		s3 := "| sort @timestamp asc"
		s4 := "| limit 10000"

		return s1 + s2 + s3 + s4
	}

	if productNumber != "" {
		s1 := "fields @timestamp, `fields.metadata.device-product-number`, `fields.metadata.device-serial-number`, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, `fields.metadata.xml-generator-object-path`"
		s2 := "| filter ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(`fields.metadata.xml-generator-object-path`) and ispresent(`fields.metadata.device-product-number`) and ispresent(`fields.metadata.device-serial-number`) and `fields.metadata.device-product-number`='{{.productNumber}}'"
		s3 := "| sort @timestamp asc"
		s4 := "| limit 10000"

		return s1 + s2 + s3 + s4
	}

	s1 := "fields @timestamp, `fields.metadata.device-product-number`, `fields.metadata.device-serial-number`, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, `fields.metadata.xml-generator-object-path`"
	s2 := "| filter ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(`fields.metadata.xml-generator-object-path`) and ispresent(`fields.metadata.device-product-number`) and ispresent(`fields.metadata.device-serial-number`)"
	s3 := "| sort @timestamp asc"
	s4 := "| limit 10000"

	return s1 + s2 + s3 + s4
}

// FetchData obtains the RTAs (in JSON format)  generated  by Cloud Connector (after processing the corresponding XML RTA file) depending on requestQueryParams.
// The method creates a variable insightsQueryParams and then calls a queryExecutor
// to perform the query. It returns the result and an error, if any.
func (rtasFetcher RtaFetcher) FetchData(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	insightsQueryParams, err := createInsightsQueryParams(requestQueryParams, rtasFetcher)
	if err != nil {
		return nil, err
	}

	result, err := rtasFetcher.queryExecutor.ExecuteQuery(insightsQueryParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}
