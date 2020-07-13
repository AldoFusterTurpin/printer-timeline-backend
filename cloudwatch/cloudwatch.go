package cloudwatch

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type QueryExecutor interface {
	ExecuteQuery(insightsQueryParams InsightsQueryParams) (*cloudwatchlogs.GetQueryResultsOutput, error)
}

type QueryExecutorImpl struct {
	svc *cloudwatchlogs.CloudWatchLogs
}

func NewQueryExecutorImpl(svc *cloudwatchlogs.CloudWatchLogs) QueryExecutorImpl {
	return QueryExecutorImpl{svc}
}

type InsightsQueryParams struct {
	StartTimeEpoch, EndTimeEpoch int64
	LogGroupName, Query          string
}

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
