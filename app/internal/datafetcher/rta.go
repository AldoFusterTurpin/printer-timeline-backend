package datafetcher

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudwatch"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/queryparams"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// RtaFetcher obtains the Rtas in JSON format (after RTA in XML has been parsed by the Cloud Connector)
// based on request query parameters.
type RtaFetcher interface {
	FetchData(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error)
	GetRtas(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error)
}

// RtaFetcherImpl is the implementation of RtaFetcher that uses a queryExecutor to perform a query
//and obtain the Rtas.
type RtaFetcherImpl struct {
	queryExecutor cloudwatch.QueryExecutor
}

// NewRtaFetcherImpl creates a new RtaFetcher implementation.
func NewRtaFetcherImpl(queryExecutor cloudwatch.QueryExecutor) RtaFetcher {
	return RtaFetcherImpl{queryExecutor}
}

// createQueryTemplate returns a new query template depending on the productNumber and serialNumber parameters.
// The resulting query template will be used by the queryExecutor to obtain the Rtas.
func (rtasFetcherImpl RtaFetcherImpl) createQueryTemplate(productNumber, serialNumber string) (queryTemplateString string) {
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

// createInsightsQueryParams creates InsightQueryParameters based on requestQueryParams.
// The returned InsightQueryParameters will be used by a QueryExecutor to execute the query. It also returns an error, if any.
func (rtasFetcherImpl RtaFetcherImpl) createInsightsQueryParams(requestQueryParams map[string]string) (insightsQueryParams cloudwatch.InsightsQueryParams, err error) {
	startTime, endTime, err := queryparams.ExtractTimeRange(requestQueryParams)
	if err != nil {
		return
	}

	productNumber, serialNumber, err := queryparams.ExtractPrinterInfo(requestQueryParams)
	if err != nil {
		return
	}

	queryTemplate := rtasFetcherImpl.createQueryTemplate(productNumber, serialNumber)

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
		LogGroupName:   "/aws/lambda/AWSUploadRTA",
		Query:          queryToExecute,
	}
	return insightsQueryParams, nil
}

// GetRtas obtains the RTAs (in JSON format)  generated  by Cloud Connector (after processing the corresponding XML RTA file) depending on requestQueryParams.
// The method creates a variable insightsQueryParams and then calls a queryExecutor
// to perform the query. It returns the result and an error, if any.
func (rtasFetcherImpl RtaFetcherImpl) GetRtas(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	insightsQueryParams, err := rtasFetcherImpl.createInsightsQueryParams(requestQueryParams)
	if err != nil {
		return nil, err
	}

	result, err := rtasFetcherImpl.queryExecutor.ExecuteQuery(insightsQueryParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rtasFetcherImpl RtaFetcherImpl) FetchData(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	return rtasFetcherImpl.GetRtas(requestQueryParams)
}
