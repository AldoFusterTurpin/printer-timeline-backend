package common

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"strconv"
)

func ConvertEpochStringToUint64(epochToConvert string) (epochConverted int64, err error) {
	return strconv.ParseInt(epochToConvert, 10, 64)
}

func CloudWatchInsightsQuery(svc *cloudwatchlogs.CloudWatchLogs, startTimeEpoch int64, endTimeEpoch int64, logGroupName string, queryString string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	startQueryInput := &cloudwatchlogs.StartQueryInput{
		StartTime:    aws.Int64(startTimeEpoch),
		EndTime:      aws.Int64(endTimeEpoch),
		LogGroupName: aws.String(logGroupName),
		QueryString:  aws.String(queryString),
	}

	startQueryOutput, err := svc.StartQuery(startQueryInput)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	queryResultsInput := &cloudwatchlogs.GetQueryResultsInput{QueryId: startQueryOutput.QueryId}
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
	fmt.Println(*queryResultsOutput.Status)
	return queryResultsOutput, nil
}
