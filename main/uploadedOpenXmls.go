package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-gonic/gin"
	"net/http"
	"text/template"
)

func prepareInsightsQueryParametersOfUploadedXmlsQuery(c *gin.Context) (startTimeEpoch int64, endTimeEpoch int64, queryString string, err error) {
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

	templateString := getQueryTemplate(productNumber, serialNumber)
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
	queryString = result.String()
	return startTimeEpoch, endTimeEpoch, queryString, err
}

func getQueryTemplate(productNumber string, serialNumber string) (templateString string) {
	if productNumber != "" && serialNumber != "" {
		return `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date) and fields.ProductNumber="{{.productNumber}}" and fields.SerialNumber="{{.serialNumber}}"
								| sort @timestamp asc
								| limit 10000`
	}
	if productNumber != "" {
		return `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date) and fields.ProductNumber="{{.productNumber}}"
								| sort @timestamp asc
								| limit 10000`
	}

	return `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date
								| filter ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date)
								| sort @timestamp asc
								| limit 10000`
}

func getUploadedOpenXmls(svc *cloudwatchlogs.CloudWatchLogs) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTimeEpoch, endTimeEpoch, queryString, err := prepareInsightsQueryParametersOfUploadedXmlsQuery(c)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		logGroupName := "/aws/lambda/AWSUpload"
		queryResultsOutput, err := cloudWatchInsightsQuery(svc, startTimeEpoch, endTimeEpoch, logGroupName, queryString)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		c.JSON(http.StatusOK, queryResultsOutput)
	}
}
