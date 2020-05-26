package common

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type CloudWatchQueryExecutor interface {
	ExecuteQuery(svc *cloudwatchlogs.CloudWatchLogs) (*cloudwatchlogs.GetQueryResultsOutput, error)
}

type CloudWatchQueryExecutorImpl struct {
	startTimeEpoch, endTimeEpoch int64
	logGroupName, queryString    string
}

func (queryExecutor* CloudWatchQueryExecutorImpl) Init(startTimeEpoch, endTimeEpoch int64, logGroupName, queryString string) {
	queryExecutor.startTimeEpoch = startTimeEpoch
	queryExecutor.endTimeEpoch = endTimeEpoch
	queryExecutor.logGroupName = logGroupName
	queryExecutor.queryString = queryString
}

func (queryExecutor* CloudWatchQueryExecutorImpl) ExecuteQuery(svc *cloudwatchlogs.CloudWatchLogs) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	startQueryInput := &cloudwatchlogs.StartQueryInput{
		StartTime:    aws.Int64(queryExecutor.startTimeEpoch),
		EndTime:      aws.Int64(queryExecutor.endTimeEpoch),
		LogGroupName: aws.String(queryExecutor.logGroupName),
		QueryString:  aws.String(queryExecutor.queryString),
	}

	startQueryOutput, err := svc.StartQuery(startQueryInput)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	queryResultsInput := &cloudwatchlogs.GetQueryResultsInput{
		QueryId: startQueryOutput.QueryId,
	}
	queryResultsOutput, err := svc.GetQueryResults(queryResultsInput)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	for *queryResultsOutput.Status == cloudwatchlogs.QueryStatusRunning || *queryResultsOutput.Status == cloudwatchlogs.QueryStatusScheduled {
		fmt.Println("INFO: Waiting query to finish")
		queryResultsOutput, err = svc.GetQueryResults(queryResultsInput)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return queryResultsOutput, nil
}
