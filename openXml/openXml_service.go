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

func PrepareInsightsQueryParameters(requestQueryParameters map[string]string) (queryParams cloudwatch.InsightsQueryParams, err error) {
	startTime, endTime, err := queryParamsCtrl.ExtractTimeRange(requestQueryParameters)
	if err != nil {
		return
	}

	productNumber, serialNumber, err := queryParamsCtrl.ExtractPrinterInfo(requestQueryParameters);
	if err != nil {
		return
	}

	templateString := selectQueryTemplate(productNumber, serialNumber)
	queryTemplate, err := template.New("queryTemplate").Parse(templateString)
	if err != nil {
		return
	}

	mapValues := map[string]interface{}{
		"productNumber": productNumber,
		"serialNumber":  serialNumber,
	}

	var query bytes.Buffer
	if err = queryTemplate.Execute(&query, mapValues); err != nil {
		return
	}
	queryParams = cloudwatch.InsightsQueryParams{
		startTime,
		endTime,
		"/aws/lambda/AWSUpload",
		query.String(),
	}
	return queryParams, nil
}

func GetUploadedOpenXmls(svc *cloudwatchlogs.CloudWatchLogs, requestQueryParams map[string]string) (int, *cloudwatchlogs.GetQueryResultsOutput) {
	insightsQueryParams, err := PrepareInsightsQueryParameters(requestQueryParams)
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
