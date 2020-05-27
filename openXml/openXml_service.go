package openXml

import (
	"bitbucket.org/aldoft/printer-timeline-backend/cloudwatch"
	"bitbucket.org/aldoft/printer-timeline-backend/errors"
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"net/http"
	"strconv"
	"text/template"
	"time"
)


func selectQueryTemplate(productNumber, serialNumber string) (templateString string) {
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

func ExtractTimeRange(queryParameters map[string]string) (startTimeEpoch int64, endTimeEpoch int64, err error) {

	timeTypeStr := queryParameters["time_type"]
	if timeTypeStr == "" {
		err = errors.QueryStringMissingTimeRangeType
		return
	}

	startTimeStr := queryParameters["start_time"]
	endTimeStr := queryParameters["end_time"]
	offsetUnits := queryParameters["offset_units"]
	offsetValue := queryParameters["offset_value"]
	switch timeTypeStr {
	case "relative":
		if startTimeStr != "" {
			err = errors.QueryStringStartTimeAppears
			return
		}
		if endTimeStr != "" {
			err = errors.QueryStringEndTimeAppears
			return
		}
		if offsetUnits == "" {
			err = errors.QueryStringMissingOffsetUnits
			return
		}
		if offsetUnits != "seconds" && offsetUnits != "minutes" {
			err = errors.QueryStringUnsupportedOffsetUnits
			return
		}
		if offsetValue == "" {
			err = errors.QueryStringMissingOffsetValue
			return
		}

		var offsetValueInt int
		offsetValueInt, err = strconv.Atoi(offsetValue)
		if err != nil {
			err = errors.QueryStringUnsupportedOffsetValue
			return
		}

		if offsetUnits == "minutes" && offsetValueInt > 60 {
			err = errors.QueryStringUnsupportedOffsetValue
			return
		}
		if offsetUnits == "seconds" && offsetValueInt > 3600 {
			err = errors.QueryStringUnsupportedOffsetValue
			return
		}
		if offsetValueInt < 1 {
			err = errors.QueryStringUnsupportedOffsetValue
			return
		}

		endTimeEpoch = time.Now().Unix()

		var duration time.Duration
		if offsetUnits == "minutes" {
			duration = -1 * time.Minute * time.Duration(offsetValueInt)
		} else if offsetUnits == "seconds" {
			duration = -1 * time.Second * time.Duration(offsetValueInt)
		}
		startTimeEpoch = time.Now().Add(duration).Unix()

	case "absolute":
		if startTimeStr == "" {
			err =  errors.QueryStringMissingStartTime
			return
		}

		startTimeEpoch, err = strconv.ParseInt(startTimeStr, 10, 64)
		if err != nil {
			err = errors.QueryStringUnsupportedStartTime
			return
		}

		if endTimeStr == "" {
			err = errors.QueryStringMissingEndTime
			return
		}
		endTimeEpoch, err = strconv.ParseInt(endTimeStr, 10, 64)
		if err != nil {
			err = errors.QueryStringUnsupportedEndTime
			return
		}

		diff := time.Unix(endTimeEpoch, 0).Sub(time.Unix(startTimeEpoch, 0))
		if diff.Minutes() > 60 {
			err = errors.QueryStringTimeDifferenceTooBig
			return
		}
		if diff.Minutes() < 0 {
			err = errors.QueryStringEndTimePreviousThanStartTime
			return
		}
	default:
		err = errors.QueryStringUnsupportedTimeRangeType
		return
	}
	return startTimeEpoch, endTimeEpoch, nil
}

func ExtractPrinterInfo(queryParameters map[string]string) (productNumber string, serialNumber string, err error) {
	productNumber = queryParameters["pn"]
	productNumber = queryParameters["sn"]
	if productNumber == "" && serialNumber != "" {
		err = errors.QueryStringPnSn
		return
	}
	return productNumber, serialNumber, nil
}

func PrepareInsightsQueryParameters(requestQueryParameters map[string]string) (queryParams cloudwatch.InsightsQueryParams, err error) {
	startTime, endTime, err := ExtractTimeRange(requestQueryParameters)
	if err != nil {
		return
	}

	productNumber, serialNumber, err := ExtractPrinterInfo(requestQueryParameters);
	if err != nil {
		return
	}

	templateString := selectQueryTemplate(productNumber, serialNumber)
	queryTemplate, err := template.New("queryTemplate").Parse(templateString)
	if err != nil {
		return
	}

	mapValues := map[string]interface{}{
		"productNumber": productNumber,
		"serialNumber":  serialNumber,
	}

	var query bytes.Buffer
	if err = queryTemplate.Execute(&query, mapValues); err != nil {
		return
	}
	queryParams = cloudwatch.InsightsQueryParams{
		startTime,
		endTime,
		"/aws/lambda/AWSUpload",
		query.String(),
	}
	return queryParams, nil
}

func GetUploadedOpenXmls(svc *cloudwatchlogs.CloudWatchLogs, requestQueryParams map[string]string) (int, *cloudwatchlogs.GetQueryResultsOutput) {
	insightsQueryParams, err := PrepareInsightsQueryParameters(requestQueryParams)
	if err != nil {
		fmt.Println(err.Error())
		return http.StatusInternalServerError, nil
	}

	result, err := cloudwatch.ExecuteQuery(svc, insightsQueryParams)
	if err != nil {
		fmt.Println(err.Error())
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, result
}
