package openXml

import (
	"bitbucket.org/aldoft/printer-timeline-backend/cloudwatch"
	"bitbucket.org/aldoft/printer-timeline-backend/queryParamsCtrl"
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"net/http"
	"text/template"
)

func selectQueryTemplate(productNumber, serialNumber string) (templateString string) {
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


func createQuery(productNumber, serialNumber string) (query string, err error){
	queryTemplateString := selectQueryTemplate(productNumber, serialNumber)

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

func GetInsightsQueryParams(requestQueryParams map[string]string) (insightsQueryParams cloudwatch.InsightsQueryParams, err error) {
	startTime, endTime, err := queryParamsCtrl.ExtractTimeRange(requestQueryParams)
	if err != nil {
		return
	}

	productNumber, serialNumber, err := queryParamsCtrl.ExtractPrinterInfo(requestQueryParams)
	if err != nil {
		return
	}

	queryToExecute, err := createQuery(productNumber, serialNumber)
	insightsQueryParams = cloudwatch.InsightsQueryParams{
		StartTimeEpoch: startTime,
		EndTimeEpoch:   endTime,
		LogGroupName:   "/aws/lambda/AWSUpload",
		Query:          queryToExecute,
	}
	return insightsQueryParams, nil
}

func GetUploadedOpenXmls(svc *cloudwatchlogs.CloudWatchLogs, requestQueryParams map[string]string) (int, *cloudwatchlogs.GetQueryResultsOutput) {
	insightsQueryParams, err := GetInsightsQueryParams(requestQueryParams)
	if err != nil {
		fmt.Println(err.Error())
		return http.StatusInternalServerError, nil
	}

	result, err := cloudwatch.ExecuteQuery(svc, insightsQueryParams)
	if err != nil {
		fmt.Println(err.Error())
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, result
}
