package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"time"
)

func defaultStartTime() time.Time {
	return time.Now().Add(-time.Minute * 5)
}

func defaultEndTime() time.Time {
	return time.Now()
}

func convertEpochStringToUint64(epochToConvert string, defaultEpoch int64) (epochConverted int64, err error) {
	if epochToConvert == "" {
		return defaultEpoch, nil
	}
	return strconv.ParseInt(epochToConvert, 10, 64)
}

func cloudWatchInsightsQuery(svc *cloudwatchlogs.CloudWatchLogs, startTimeEpoch int64, endTimeEpoch int64, logGroupName string, queryString string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
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

func initRouter(svc *cloudwatchlogs.CloudWatchLogs) *gin.Engine {
	router := gin.Default()
	router.GET("api/open_xml", getUploadedOpenXmls(svc))
	return router
}

func main() {
	awsRegion, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		fmt.Println(errors.New("could not load BUCKET_REGION env var"))
		return
	}

	var err error
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	svc := cloudwatchlogs.New(sess)

	router := initRouter(svc)

	if err = router.Run(); err != nil {
		fmt.Println(err.Error())
		return
	}
}