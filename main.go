package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"text/template"
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
		startTime := c.Query("start_time")
		startTimeEpoch, startTimeError := convertEpochStringToUint64(startTime, defaultStartTime().Unix())
		if startTimeError != nil {
			fmt.Println(startTimeError.Error())
			return
		}

		endTime := c.Query("end_time")
		endTimeEpoch, endTimeError := convertEpochStringToUint64(endTime, defaultEndTime().Unix())
		if endTimeError != nil {
			fmt.Println(endTimeError.Error())
			return
		}

		productNumber := c.Query("pn")
		serialNumber := c.Query("sn")

		var templateString string
		if productNumber != "" && serialNumber != "" {
			templateString = `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date) and fields.ProductNumber="{{.productNumber}}" and fields.SerialNumber="{{.serialNumber}}"
								| sort @timestamp asc
								| limit 10000`
		} else if  productNumber != "" {
			templateString = `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date) and fields.ProductNumber="{{.productNumber}}"
								| sort @timestamp asc
								| limit 10000`
		} else {
			templateString = `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date)
								| sort @timestamp asc
								| limit 10000`
		}
		myTemplate, err := template.New("myTemplate").Parse(templateString)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		templateMapValues := map[string]interface{}{
			"productNumber": productNumber,
			"serialNumber":  serialNumber,
		}
		var result bytes.Buffer
		if err := myTemplate.Execute(&result, templateMapValues); err != nil {
			fmt.Println(err.Error())
		}
		queryString := result.String()
		logGroupName := "/aws/lambda/AWSUpload"
		queryResultsOutput, err := cloudWatchInsightsQuery(svc, startTimeEpoch, endTimeEpoch, logGroupName, queryString)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		c.JSON(http.StatusOK, queryResultsOutput)
	}
}

func cloudWatchInsightsQuery(svc *cloudwatchlogs.CloudWatchLogs, startTimeEpoch int64, endTimeEpoch int64, logGroupName string, queryString string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	startQueryInput := &cloudwatchlogs.StartQueryInput{
		StartTime:    aws.Int64(startTimeEpoch),
		EndTime:      aws.Int64(endTimeEpoch),
		LogGroupName: aws.String(logGroupName),
		QueryString:  aws.String(queryString),
	}

	startQueryOutput, startQueryError := svc.StartQuery(startQueryInput)
	if startQueryError != nil {
		fmt.Println(startQueryError.Error())
		return nil, startQueryError
	}

	queryResultsInput := &cloudwatchlogs.GetQueryResultsInput{QueryId: startQueryOutput.QueryId}
	queryResultsOutput, getQueryResultsError := svc.GetQueryResults(queryResultsInput)
	if getQueryResultsError != nil {
		fmt.Println(getQueryResultsError.Error())
		return nil, getQueryResultsError
	}
	for *queryResultsOutput.Status == cloudwatchlogs.QueryStatusRunning || *queryResultsOutput.Status == cloudwatchlogs.QueryStatusScheduled {
		fmt.Println("INFO: Waiting query to finish")
		queryResultsOutput, getQueryResultsError = svc.GetQueryResults(queryResultsInput)
		if getQueryResultsError != nil {
			fmt.Println(getQueryResultsError.Error())
			return nil, getQueryResultsError
		}
	}
	fmt.Println(*queryResultsOutput.Status)
	return queryResultsOutput, nil
}

func setUpRouter(svc *cloudwatchlogs.CloudWatchLogs) *gin.Engine {
	router := gin.Default()
	router.GET("api/open_xml", queryUploadedOpenXmls(svc))
	return router
}

func main() {
	sess, sessionError := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if sessionError != nil {
		fmt.Println(sessionError.Error())
	}

	svc := cloudwatchlogs.New(sess)

	router := setUpRouter(svc)
	routerError := router.Run()
	if routerError != nil {
		fmt.Println(routerError.Error())
	}
}
