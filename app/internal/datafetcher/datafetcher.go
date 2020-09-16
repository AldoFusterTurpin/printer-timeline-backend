// Package datafetcher provides an interface to obtain data needed by other parts of the application
package datafetcher

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudwatch"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/configs"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/queryparams"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// DataFetcher is an interface responsible of obtaining the data. Different structs will have a different logic to
// obtain its data based in the concrete implementation.
type DataFetcher interface {
	FetchData(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error)
	CreateQueryTemplate(productNumber string, serialNumber string) (queryTemplateStr string)
	GetLogGroupName() (logGroupName string)
}

// createInsightsQueryParams creates InsightQueryParameters based on requestQueryParams and the dataFetcher parameter.
// The returned InsightQueryParameters will be used by a QueryExecutor to execute the query. It also returns an error, if any.
func createInsightsQueryParams(requestQueryParams map[string]string, dataFetcher DataFetcher) (insightsQueryParams cloudwatch.InsightsQueryParams, err error) {
	startTime, endTime, err := queryparams.ExtractTimeRange(requestQueryParams, configs.GetMaxTimeDiffInMinutes())
	if err != nil {
		return
	}

	productNumber, serialNumber, err := queryparams.ExtractPrinterInfo(requestQueryParams)
	if err != nil {
		return
	}

	mapValues := map[string]string{
		"productNumber": productNumber,
		"serialNumber":  serialNumber,
	}

	queryTemplate := dataFetcher.CreateQueryTemplate(productNumber, serialNumber)

	queryToExecute, err := cloudwatch.CreateQuery(queryTemplate, mapValues)
	if err != nil {
		return
	}

	insightsQueryParams = cloudwatch.InsightsQueryParams{
		StartTimeEpoch: startTime.Unix(),
		EndTimeEpoch:   endTime.Unix(),
		LogGroupName:   dataFetcher.GetLogGroupName(),
		Query:          queryToExecute,
	}
	return insightsQueryParams, nil
}
