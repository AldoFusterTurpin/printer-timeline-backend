package cloudwatch

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type InsightsQueryParams struct {
	StartTimeEpoch, EndTimeEpoch int64
	LogGroupName, QueryString    string
}

func ExecuteQuery(svc *cloudwatchlogs.CloudWatchLogs, queryParams InsightsQueryParams) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	startQueryInput := &cloudwatchlogs.StartQueryInput{
		StartTime:    aws.Int64(queryParams.StartTimeEpoch),
		EndTime:      aws.Int64(queryParams.EndTimeEpoch),
		LogGroupName: aws.String(queryParams.LogGroupName),
		QueryString:  aws.String(queryParams.QueryString),
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
