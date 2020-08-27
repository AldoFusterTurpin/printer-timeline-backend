// Package cloudJson provides an interface to obtain the CloudJsons based on request parameters.
// It will handle the requests of Jsons created by the Cloud Connector.
package cloudJson

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudwatch"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/queryparams"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// CloudJsonsFetcher obtains the Cloud Jsons created by the Cloud Connector based on request query parameters.
type CloudJsonsFetcher interface {
	FetchData(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error)
	GetCloudJsons(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error)
}

// CloudJsonsFetcherImpl is the implementation of CloudJsonsFetcher that uses a queryExecutor to perform a query
//and obtain the Cloud Jsons.
type CloudJsonsFetcherImpl struct {
	queryExecutor cloudwatch.QueryExecutor
}

// NewCloudJsonsFetcherImpl creates a new CloudJsonsFetcher implementationm
func NewCloudJsonsFetcherImpl(queryExecutor cloudwatch.QueryExecutor) CloudJsonsFetcher {
	return CloudJsonsFetcherImpl{queryExecutor}
}

// createQueryTemplate returns a new query template depending on the productNumber and serialNumber parameters.
// The resulting query template will be used by the queryExecutor to obtain the Cloud Jsons.
func (cloudJsonsFetcherImpl CloudJsonsFetcherImpl) createQueryTemplate(productNumber, serialNumber string) (queryTemplateString string) {
	// As the query string contains backticks ``, I need to surround those backticks with double quotes("").
	// There is also the expression "{{.productNumber}}" which need to be
	// surrounded by double quotes ("") (because the Go backticks used to create a string with multiple lines can't contain double quotes inside it).
	// Last, I need to concatenate all those strings but need to create temporary variables because are 'untyped string constants'.
	// For more info, check: https://blog.golang.org/constants
	if productNumber != "" && serialNumber != "" {
		s1 := `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date, `
		s2 := "`fields.metadata.xml-generator-object-path`"
		s3 := `| filter (ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date)`
		s4 := " and ispresent(`fields.metadata.xml-generator-object-path`)"
		s5 := ` and fields.topic = "json" and fields.ProductNumber="{{.productNumber}}" and fields.SerialNumber="{{.serialNumber}}")
		| sort @timestamp asc
		| limit 10000`

		return s1 + s2 + s3 + s4 + s5
	}
	if productNumber != "" {
		s1 := `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date, `
		s2 := "`fields.metadata.xml-generator-object-path`"
		s3 := `| filter (ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date)`
		s4 := " and ispresent(`fields.metadata.xml-generator-object-path`)"
		s5 := ` and fields.topic = "json" and fields.ProductNumber="{{.productNumber}}")
		| sort @timestamp asc
		| limit 10000`

		return s1 + s2 + s3 + s4 + s5
	}

	s1 := `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date, `
	s2 := "`fields.metadata.xml-generator-object-path`"
	s3 := `| filter (ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date)`
	s4 := " and ispresent(`fields.metadata.xml-generator-object-path`)"
	s5 := ` and fields.topic = "json")`
	s6 := `| sort @timestamp asc | limit 10000`

	return s1 + s2 + s3 + s4 + s5 + s6
}

// createInsightsQueryParams creates InsightQueryParameters based on requestQueryParams.
// The returned InsightQueryParameters will be used by a QueryExecutor to execute the query. It also returns an error, if any.
func (cloudJsonsFetcherImpl CloudJsonsFetcherImpl) createInsightsQueryParams(requestQueryParams map[string]string) (insightsQueryParams cloudwatch.InsightsQueryParams, err error) {
	startTime, endTime, err := queryparams.ExtractTimeRange(requestQueryParams)
	if err != nil {
		return
	}

	productNumber, serialNumber, err := queryparams.ExtractPrinterInfo(requestQueryParams)
	if err != nil {
		return
	}

	queryTemplate := cloudJsonsFetcherImpl.createQueryTemplate(productNumber, serialNumber)

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
		LogGroupName:   "/aws/lambda/AWSParser",
		Query:          queryToExecute,
	}
	return insightsQueryParams, nil
}

// GetCloudJsons obtains the Jsons created by Cloud Connector depending on requestQueryParams.
// The method basically creates a variable insightsQueryParams and then calls a queryExecutor
// to perform the query. It returns the result and an error, if any.
func (cloudJsonsFetcherImpl CloudJsonsFetcherImpl) GetCloudJsons(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	insightsQueryParams, err := cloudJsonsFetcherImpl.createInsightsQueryParams(requestQueryParams)
	if err != nil {
		return nil, err
	}

	result, err := cloudJsonsFetcherImpl.queryExecutor.ExecuteQuery(insightsQueryParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cloudJsonsFetcherImpl CloudJsonsFetcherImpl) FetchData(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	return cloudJsonsFetcherImpl.GetCloudJsons(requestQueryParams)
}
