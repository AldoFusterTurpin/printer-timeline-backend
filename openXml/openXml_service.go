package openXml

import (
	"bytes"
	"text/template"

	"bitbucket.org/aldoft/printer-timeline-backend/cloudwatch"
	"bitbucket.org/aldoft/printer-timeline-backend/queryParamsCtrl"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type OpenXmlsFetcher interface {
	GetUploadedOpenXmls(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error)
}

type OpenXmlsFetcherImpl struct {
	queryExecutor cloudwatch.QueryExecutor
}

func NewOpenXmlsFetcherImpl(queryExecutor cloudwatch.QueryExecutor) OpenXmlsFetcherImpl {
	return OpenXmlsFetcherImpl{queryExecutor}
}

func (openXmlsFetcherImpl OpenXmlsFetcherImpl) selectQueryTemplate(productNumber, serialNumber string) (templateString string) {
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

func (openXmlsFetcherImpl OpenXmlsFetcherImpl) createQuery(productNumber, serialNumber string) (query string, err error) {
	queryTemplateString := openXmlsFetcherImpl.selectQueryTemplate(productNumber, serialNumber)

	queryTemplate, err := template.New("queryTemplate").Parse(queryTemplateString)
	if err != nil {
		return
	}

	mapValues := map[string]string{
		"productNumber": productNumber,
		"serialNumber":  serialNumber,
	}

	var queryBuffer bytes.Buffer
	err = queryTemplate.Execute(&queryBuffer, mapValues)
	if err != nil {
		return
	}
	return queryBuffer.String(), nil
}

func (openXmlsFetcherImpl OpenXmlsFetcherImpl) getInsightsQueryParams(requestQueryParams map[string]string) (insightsQueryParams cloudwatch.InsightsQueryParams, err error) {
	startTime, endTime, err := queryParamsCtrl.ExtractTimeRange(requestQueryParams)
	if err != nil {
		return
	}

	productNumber, serialNumber, err := queryParamsCtrl.ExtractPrinterInfo(requestQueryParams)
	if err != nil {
		return
	}

	queryToExecute, err := openXmlsFetcherImpl.createQuery(productNumber, serialNumber)
	insightsQueryParams = cloudwatch.InsightsQueryParams{
		StartTimeEpoch: startTime.Unix(),
		EndTimeEpoch:   endTime.Unix(),
		LogGroupName:   "/aws/lambda/AWSUpload",
		Query:          queryToExecute,
	}
	return insightsQueryParams, nil
}

func (openXmlsFetcherImpl OpenXmlsFetcherImpl) GetUploadedOpenXmls(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	insightsQueryParams, err := openXmlsFetcherImpl.getInsightsQueryParams(requestQueryParams)
	if err != nil {
		return nil, err
	}

	result, err := openXmlsFetcherImpl.queryExecutor.ExecuteQuery(insightsQueryParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}
