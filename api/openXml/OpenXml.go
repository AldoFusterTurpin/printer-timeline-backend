package openXml

import (
	"bitbucket.org/aldoft/printer-timeline-backend/api/common"
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-gonic/gin"
	"net/http"
	"text/template"
)


type QueryPrinterInfo struct {
	StartTimeEpoch int64
	EndTimeEpoch int64
	ProductNumber string
	SerialNumber string
}

func selectQueryTemplate(productNumber string, serialNumber string) (templateString string) {
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

func getInfoFromQueryStrings(c *gin.Context) (*QueryPrinterInfo, error) {
	queryPrinterInfo := new(QueryPrinterInfo)

	timeTypeStr := c.Query("time_type")
	if timeTypeStr == "" {
		return nil, common.QueryStringMissingTimeRangeType
	}

	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	switch timeTypeStr {
	case "relative":
		if startTimeStr == "" {
			queryPrinterInfo.StartTimeEpoch = common.DefaultStartTime().Unix()
		} else {
			var err error
			queryPrinterInfo.StartTimeEpoch, err = common.ConvertEpochStringToUint64(startTimeStr)
			if err != nil {
				return nil, err
			}
		}

		if endTimeStr != "" {
			return nil, common.QueryStringPresentEndTime
		}
		queryPrinterInfo.EndTimeEpoch = common.DefaultEndTime().Unix()

	case "absolute":
		if startTimeStr == "" {
			return nil, common.QueryStringMissingStartTime
		}
		var err error
		queryPrinterInfo.StartTimeEpoch, err = common.ConvertEpochStringToUint64(startTimeStr)
		if err != nil {
			return nil, err
		}

		if endTimeStr == "" {
			return nil, common.QueryStringMissingEndTime
		}
		queryPrinterInfo.EndTimeEpoch, err = common.ConvertEpochStringToUint64(endTimeStr)
		if err != nil {
			return nil, err
		}
	default:
		return nil, common.QueryStringUnsupportedTimeRangeType
	}

	queryPrinterInfo.ProductNumber = c.Query("pn")
	queryPrinterInfo.SerialNumber = c.Query("sn")
	if queryPrinterInfo.ProductNumber == "" && queryPrinterInfo.SerialNumber != ""{
		return nil, common.QueryStringPnSn
	}

	return queryPrinterInfo, nil
}


func prepareInsightsQueryParameters(c *gin.Context) (startTimeEpoch int64, endTimeEpoch int64, awsInsightsQuery string, err error) {
	queryPrinterInfo, err := getInfoFromQueryStrings(c)
	if err != nil {
		return
	}

	templateString := selectQueryTemplate(queryPrinterInfo.ProductNumber, queryPrinterInfo.SerialNumber)
	queryTemplate, err := template.New("queryTemplate").Parse(templateString)
	if err != nil {
		return
	}

	mapValues := map[string]interface{} {
		"productNumber": queryPrinterInfo.ProductNumber,
		"serialNumber":  queryPrinterInfo.SerialNumber,
	}

	var query bytes.Buffer
	if err = queryTemplate.Execute(&query, mapValues); err != nil {
		return
	}

	return queryPrinterInfo.StartTimeEpoch, queryPrinterInfo.EndTimeEpoch, query.String(), nil
}


func GetUploadedOpenXmls(svc *cloudwatchlogs.CloudWatchLogs) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTimeEpoch, endTimeEpoch, queryString, err := prepareInsightsQueryParameters(c)
		if err != nil {
			fmt.Println(err.Error())
			jsonResponse := gin.H{
				"error": err.Error(),
			}
			c.JSON(http.StatusInternalServerError,  jsonResponse)
			return
		}
		logGroupName := "/aws/lambda/AWSUpload"
		queryResultsOutput, err := common.CloudWatchInsightsQuery(svc, startTimeEpoch, endTimeEpoch, logGroupName, queryString)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		c.JSON(http.StatusOK, queryResultsOutput)
	}
}
