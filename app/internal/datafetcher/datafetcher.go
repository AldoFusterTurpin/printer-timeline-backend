// Package datafetcher provides an interface to obtain data needed by other parts of the application
package datafetcher

import "github.com/aws/aws-sdk-go/service/cloudwatchlogs"

type DataFetcher interface {
	FetchData(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error)
}
