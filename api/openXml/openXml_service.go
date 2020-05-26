package openXml

import (
	"bitbucket.org/aldoft/printer-timeline-backend/api/common"
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"net/http"
	"strconv"
	"text/template"
	"time"
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

func getInfoFromQueryStrings(queryParameters map[string]string) (*QueryPrinterInfo, error) {
	queryPrinterInfo := new(QueryPrinterInfo)

	timeTypeStr := queryParameters["time_type"]
	if timeTypeStr == "" {
		return nil, common.QueryStringMissingTimeRangeTypeError
	}

	startTimeStr := queryParameters["start_time"]
	endTimeStr := queryParameters["end_time"]
	offsetUnits := queryParameters["offset_units"]
	offsetValue := queryParameters["offset_value"]
	switch timeTypeStr {
	case "relative":
		if startTimeStr != "" {
			return nil, common.QueryStringStartTimeAppearsError
		}
		if endTimeStr != "" {
			return nil, common.QueryStringEndTimeAppearsError
		}
		if offsetUnits == "" {
			return nil, common.QueryStringMissingOffsetUnitsError
		}
		if offsetUnits != "seconds" && offsetUnits != "minutes" {
			return nil, common.QueryStringUnsupportedOffsetUnitsError
		}
		if offsetValue == "" {
			return nil, common.QueryStringMissingOffsetValueError
		}

		offsetValueInt, err := strconv.Atoi(offsetValue)
		if  err != nil {
			return nil, common.QueryStringUnsupportedOffsetValueError
		}

		if offsetUnits == "minutes" && offsetValueInt > 60 {
			return nil, common.QueryStringUnsupportedOffsetValueError
		}
		if offsetUnits == "seconds" && offsetValueInt > 3600 {
			return nil, common.QueryStringUnsupportedOffsetValueError
		}
		if offsetValueInt < 1 {
			return nil, common.QueryStringUnsupportedOffsetValueError
		}

		queryPrinterInfo.EndTimeEpoch = time.Now().Unix()

		var duration time.Duration
		if offsetUnits == "minutes" {
			duration = -1 * time.Minute * time.Duration(offsetValueInt)
		} else if offsetUnits == "seconds" {
			duration = -1 * time.Second * time.Duration(offsetValueInt)
		} else {
			return nil, common.QueryStringUnsupportedOffsetUnitsError
		}
		queryPrinterInfo.StartTimeEpoch = time.Now().Add(duration).Unix()


	case "absolute":
		if startTimeStr == "" {
			return nil, common.QueryStringMissingStartTimeError
		}

		var err error
		queryPrinterInfo.StartTimeEpoch, err = common.ConvertEpochStringToUint64(startTimeStr)
		if err != nil {
			return nil, common.QueryStringUnsupportedStartTimeError
		}

		if endTimeStr == "" {
			return nil, common.QueryStringMissingEndTimeError
		}
		queryPrinterInfo.EndTimeEpoch, err = common.ConvertEpochStringToUint64(endTimeStr)
		if err != nil {
			return nil, common.QueryStringUnsupportedEndTimeError
		}

		diff := time.Unix(queryPrinterInfo.EndTimeEpoch, 0).Sub(time.Unix(queryPrinterInfo.StartTimeEpoch, 0))
		if diff.Minutes() > 60 {
			return nil, common.QueryStringTimeDifferenceTooBig
		}
		if diff.Minutes() < 0 {
			return nil, common.QueryStringEndTimePreviousThanStartTime
		}
	default:
		return nil, common.QueryStringUnsupportedTimeRangeTypeError
	}

	queryPrinterInfo.ProductNumber = queryParameters["pn"]
	queryPrinterInfo.SerialNumber = queryParameters["sn"]
	if queryPrinterInfo.ProductNumber == "" && queryPrinterInfo.SerialNumber != ""{
		return nil, common.QueryStringPnSnError
	}

	return queryPrinterInfo, nil
}


func PrepareInsightsQueryParameters(requestQueryParameters map[string]string) (startTimeEpoch int64, endTimeEpoch int64, awsInsightsQuery string, err error) {
	queryPrinterInfo, err := getInfoFromQueryStrings(requestQueryParameters)
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


/*func Handler(svc *cloudwatchlogs.CloudWatchLogs) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParameters := map[string]string {
			"time_type" : c.Query("time_type"),
			"offset_units" : c.Query("offset_units"),
			"offset_value" : c.Query("offset_value"),
			"start_time" : c.Query("start_time"),
			"end_time" : c.Query("end_time"),
			"pn" : c.Query("pn"),
			"sn" : c.Query("sn"),
		}

		startTimeEpoch, endTimeEpoch, queryString, err := PrepareInsightsQueryParameters(queryParameters)
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
}*/

func Handler(svc *cloudwatchlogs.CloudWatchLogs, queryParameters map[string]string) (int, *cloudwatchlogs.GetQueryResultsOutput) {
		startTimeEpoch, endTimeEpoch, queryString, err := PrepareInsightsQueryParameters(queryParameters)
		if err != nil {
			fmt.Println(err.Error())
			return http.StatusInternalServerError, nil
		}

		logGroupName := "/aws/lambda/AWSUpload"
		queryResultsOutput, err := common.CloudWatchInsightsQuery(svc, startTimeEpoch, endTimeEpoch, logGroupName, queryString)
		if err != nil {
			fmt.Println(err.Error())
			return http.StatusInternalServerError, nil
		}
		return http.StatusOK, queryResultsOutput
}
