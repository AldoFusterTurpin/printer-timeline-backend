package cloudwatch

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type InsightsQueryParams struct {
	StartTimeEpoch, EndTimeEpoch int64
	LogGroupName, Query          string
}

func ExecuteQuery(svc *cloudwatchlogs.CloudWatchLogs, insightsQueryParams InsightsQueryParams) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	startQueryInput := &cloudwatchlogs.StartQueryInput{
		StartTime:    aws.Int64(insightsQueryParams.StartTimeEpoch),
		EndTime:      aws.Int64(insightsQueryParams.EndTimeEpoch),
		LogGroupName: aws.String(insightsQueryParams.LogGroupName),
		QueryString:  aws.String(insightsQueryParams.Query),
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
