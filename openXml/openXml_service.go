// Package openXml provides an interface to obtain the UploadedOpenXmls based on request parameters.
// It will handle the requests of Uploaded OpenXmls.
package openXml

import (
	"bitbucket.org/aldoft/printer-timeline-backend/cloudwatch"
	"bitbucket.org/aldoft/printer-timeline-backend/queryParamsCtrl"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// OpenXmlsFetcher obtains the Uploaded OpenXMls based on request query parameters.
type OpenXmlsFetcher interface {
	GetUploadedOpenXmls(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error)
}

// OpenXmlsFetcherImpl is the implementation of OpenXmlsFetcher that uses a queryExecutor to perform a query
//and obtain the OpenXMls.
type OpenXmlsFetcherImpl struct {
	queryExecutor cloudwatch.QueryExecutor
}

// NewOpenXmlsFetcherImpl creates a new OpenXmlsFetcher implementationm
func NewOpenXmlsFetcherImpl(queryExecutor cloudwatch.QueryExecutor) OpenXmlsFetcher {
	return OpenXmlsFetcherImpl{queryExecutor}
}

// createQueryTemplate returns a new query template depending on the productNumber and serialNumber parameters.
// The resulting query template will be used by the queryExecutor to obtain the OpenXMls.
func (openXmlsFetcherImpl OpenXmlsFetcherImpl) createQueryTemplate(productNumber, serialNumber string) (queryTemplateString string) {
	if productNumber != "" && serialNumber != "" {
		return `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date) and fields.ProductNumber="{{.productNumber}}" and fields.SerialNumber="{{.serialNumber}}"
								| sort @timestamp asc
								| limit 10000`
	}
	if productNumber != "" {
		return `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date) and fields.ProductNumber="{{.productNumber}}"
								| sort @timestamp asc
								| limit 10000`
	}

	return `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date)
								| sort @timestamp asc
								| limit 10000`
}

// createInsightsQueryParams creates InsightQueryParameters based on requestQueryParams.
// The returned InsightQueryParameters will be used by a QueryExecutor to execute the query. It also returns an error, if any.
func (openXmlsFetcherImpl OpenXmlsFetcherImpl) createInsightsQueryParams(requestQueryParams map[string]string) (insightsQueryParams cloudwatch.InsightsQueryParams, err error) {
	startTime, endTime, err := queryParamsCtrl.ExtractTimeRange(requestQueryParams)
	if err != nil {
		return
	}

	productNumber, serialNumber, err := queryParamsCtrl.ExtractPrinterInfo(requestQueryParams)
	if err != nil {
		return
	}

	queryTemplate := openXmlsFetcherImpl.createQueryTemplate(productNumber, serialNumber)

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

// GetUploadedOpenXmls obtains the uploaded OpenXml depending on requestQueryParams.
// The method basically creates a variable insightsQueryParams and then calls a queryExecutor
// to perform the query. It returns the result and an error, if any.
func (openXmlsFetcherImpl OpenXmlsFetcherImpl) GetUploadedOpenXmls(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	insightsQueryParams, err := openXmlsFetcherImpl.createInsightsQueryParams(requestQueryParams)
	if err != nil {
		return nil, err
	}

	result, err := openXmlsFetcherImpl.queryExecutor.ExecuteQuery(insightsQueryParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}
