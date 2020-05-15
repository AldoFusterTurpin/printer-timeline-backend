package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func convertEpochStringToUint64(epochToConvert string, defaultEpoch int64) (epochConverted int64, err error) {
	if epochToConvert == "" {
		return defaultEpoch, nil
	}
	return strconv.ParseInt(epochToConvert, 10, 64)
}


func defaultStartTime() time.Time {
	return time.Now().Add(-time.Minute * 15)
}


func defaultEndTime() time.Time {
	return time.Now()
}


func queryUploadedOpenXmls(svc *cloudwatchlogs.CloudWatchLogs) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := c.Query("startTime")
		startTimeEpoch, startTimeError := convertEpochStringToUint64(startTime, defaultStartTime().Unix())
		 if startTimeError != nil {
			fmt.Println(startTimeError.Error())
			return
		}

		endTime := c.Query("endTime")
		endTimeEpoch, endTimeError := convertEpochStringToUint64(endTime, defaultEndTime().Unix())
		if endTimeError != nil {
			fmt.Println(endTimeError.Error())
			return
		}

		logGroupName := "/aws/lambda/AWSUpload"
		queryString := `fields @timestamp, @message | sort @timestamp desc | limit 20` //TODO change query
		startQueryInput := &cloudwatchlogs.StartQueryInput {
			StartTime: aws.Int64(startTimeEpoch),
			EndTime: aws.Int64(endTimeEpoch),
			LogGroupName: aws.String(logGroupName),
			QueryString: aws.String(queryString),
		}

		startQueryOutput, startQueryError := svc.StartQuery(startQueryInput)
		if startQueryError != nil {
			fmt.Println(startQueryError.Error())
			return
		}

		queryResultsInput := &cloudwatchlogs.GetQueryResultsInput{QueryId: startQueryOutput.QueryId}
		queryResultsOutput, getQueryResultsError := svc.GetQueryResults(queryResultsInput)
		if getQueryResultsError != nil {
			fmt.Println(getQueryResultsError.Error())
			return
		}
		for *queryResultsOutput.Status == cloudwatchlogs.QueryStatusRunning || *queryResultsOutput.Status == cloudwatchlogs.QueryStatusScheduled {
			fmt.Print("Inside waiting")
			queryResultsOutput, getQueryResultsError = svc.GetQueryResults(queryResultsInput)
			if getQueryResultsError != nil {
				fmt.Println(getQueryResultsError.Error())
				return
			}
		}
		fmt.Println(queryResultsOutput)

		c.String(http.StatusOK, "Hi")
	}
}

func main() {
	sess, sessionError := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if sessionError != nil {
		fmt.Println(sessionError.Error())
	}
	
	svc := cloudwatchlogs.New(sess)

	router := gin.Default()

	router.GET("api/open_xml", queryUploadedOpenXmls(svc))

	if routerError := router.Run(); routerError != nil {
		fmt.Println(routerError.Error())
	}
}