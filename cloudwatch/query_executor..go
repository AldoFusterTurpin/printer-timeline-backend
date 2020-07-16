// Package cloudwatch defines the functionality to perform queries in an external Service.
// It also contains a default implementation that uses AWS CloudWatch client to perform the queries in AWS Insights.
package cloudwatch

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// InsightsQueryParams are the parameters needed by QueryExecutor to execute the corresponding query
type InsightsQueryParams struct {
	StartTimeEpoch, EndTimeEpoch int64
	LogGroupName, Query          string
}

// QueryExecutor is an interface responsible of performing a query to obtain a results based on insightsQueryParams
type QueryExecutor interface {
	ExecuteQuery(insightsQueryParams InsightsQueryParams) (*cloudwatchlogs.GetQueryResultsOutput, error)
}

// QueryExecutorImpl is the default implementation of the interface that obtains the result
// using AWS CloudWatch Insights service client (svc variable) (https://docs.aws.amazon.com/sdk-for-go/api/service/cloudwatch/)
type QueryExecutorImpl struct {
	svc *cloudwatchlogs.CloudWatchLogs
}

// NewQueryExecutorImpl creates a new QueryExecutorImpl variable
func NewQueryExecutorImpl(svc *cloudwatchlogs.CloudWatchLogs) QueryExecutor {
	return QueryExecutorImpl{svc}
}

// ExecuteQuery method from QueryExecutorImpl executes a query using its cloudwatchlogs service client based on insightsQueryParams
func (queryExecutor QueryExecutorImpl) ExecuteQuery(insightsQueryParams InsightsQueryParams) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	startQueryInput := &cloudwatchlogs.StartQueryInput{
		StartTime:    aws.Int64(insightsQueryParams.StartTimeEpoch),
		EndTime:      aws.Int64(insightsQueryParams.EndTimeEpoch),
		LogGroupName: aws.String(insightsQueryParams.LogGroupName),
		QueryString:  aws.String(insightsQueryParams.Query),
	}

	startQueryOutput, err := queryExecutor.svc.StartQuery(startQueryInput)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	queryResultsInput := &cloudwatchlogs.GetQueryResultsInput{
		QueryId: startQueryOutput.QueryId,
	}
	queryResultsOutput, err := queryExecutor.svc.GetQueryResults(queryResultsInput)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	for *queryResultsOutput.Status == cloudwatchlogs.QueryStatusRunning || *queryResultsOutput.Status == cloudwatchlogs.QueryStatusScheduled {
		fmt.Println("INFO: Waiting query to finish")
		queryResultsOutput, err = queryExecutor.svc.GetQueryResults(queryResultsInput)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return queryResultsOutput, nil
}
