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

type Service interface {
	GetUploadedOpenXmls(svc *cloudwatchlogs.CloudWatchLogs, queryParameters map[string]string) (resultStatus int, resultData *cloudwatchlogs.GetQueryResultsOutput)
}

type ServiceImpl struct {
	StartTimeEpoch, EndTimeEpoch                           int64
	ProductNumber, SerialNumber, QueryString, LogGroupName string
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

func (openXmlService *ServiceImpl) Init(queryParameters map[string]string) error {

	timeTypeStr := queryParameters["time_type"]
	if timeTypeStr == "" {
		return errors.QueryStringMissingTimeRangeTypeError
	}

	startTimeStr := queryParameters["start_time"]
	endTimeStr := queryParameters["end_time"]
	offsetUnits := queryParameters["offset_units"]
	offsetValue := queryParameters["offset_value"]
	switch timeTypeStr {
	case "relative":
		if startTimeStr != "" {
			return errors.QueryStringStartTimeAppearsError
		}
		if endTimeStr != "" {
			return errors.QueryStringEndTimeAppearsError
		}
		if offsetUnits == "" {
			return errors.QueryStringMissingOffsetUnitsError
		}
		if offsetUnits != "seconds" && offsetUnits != "minutes" {
			return errors.QueryStringUnsupportedOffsetUnitsError
		}
		if offsetValue == "" {
			return errors.QueryStringMissingOffsetValueError
		}

		offsetValueInt, err := strconv.Atoi(offsetValue)
		if err != nil {
			return errors.QueryStringUnsupportedOffsetValueError
		}

		if offsetUnits == "minutes" && offsetValueInt > 60 {
			return errors.QueryStringUnsupportedOffsetValueError
		}
		if offsetUnits == "seconds" && offsetValueInt > 3600 {
			return errors.QueryStringUnsupportedOffsetValueError
		}
		if offsetValueInt < 1 {
			return errors.QueryStringUnsupportedOffsetValueError
		}

		openXmlService.EndTimeEpoch = time.Now().Unix()

		var duration time.Duration
		if offsetUnits == "minutes" {
			duration = -1 * time.Minute * time.Duration(offsetValueInt)
		} else if offsetUnits == "seconds" {
			duration = -1 * time.Second * time.Duration(offsetValueInt)
		} else {
			return errors.QueryStringUnsupportedOffsetUnitsError
		}
		openXmlService.StartTimeEpoch = time.Now().Add(duration).Unix()

	case "absolute":
		if startTimeStr == "" {
			return errors.QueryStringMissingStartTimeError
		}

		var err error
		openXmlService.StartTimeEpoch, err = strconv.ParseInt(startTimeStr, 10, 64)
		if err != nil {
			return errors.QueryStringUnsupportedStartTimeError
		}

		if endTimeStr == "" {
			return errors.QueryStringMissingEndTimeError
		}
		openXmlService.EndTimeEpoch, err = strconv.ParseInt(endTimeStr, 10, 64)
		if err != nil {
			return errors.QueryStringUnsupportedEndTimeError
		}

		diff := time.Unix(openXmlService.EndTimeEpoch, 0).Sub(time.Unix(openXmlService.StartTimeEpoch, 0))
		if diff.Minutes() > 60 {
			return errors.QueryStringTimeDifferenceTooBig
		}
		if diff.Minutes() < 0 {
			return errors.QueryStringEndTimePreviousThanStartTime
		}
	default:
		return errors.QueryStringUnsupportedTimeRangeTypeError
	}

	openXmlService.ProductNumber = queryParameters["pn"]
	openXmlService.SerialNumber = queryParameters["sn"]
	if openXmlService.ProductNumber == "" && openXmlService.SerialNumber != "" {
		return errors.QueryStringPnSnError
	}

	return nil
}

func (openXmlService *ServiceImpl) PrepareInsightsQueryParameters(requestQueryParameters map[string]string) (err error) {
	if err = openXmlService.Init(requestQueryParameters); err != nil {
		return
	}

	templateString := selectQueryTemplate(openXmlService.ProductNumber, openXmlService.SerialNumber)
	queryTemplate, err := template.New("queryTemplate").Parse(templateString)
	if err != nil {
		return
	}

	mapValues := map[string]interface{}{
		"productNumber": openXmlService.ProductNumber,
		"serialNumber":  openXmlService.SerialNumber,
	}

	var query bytes.Buffer
	if err = queryTemplate.Execute(&query, mapValues); err != nil {
		return
	}
	openXmlService.QueryString = query.String()
	return nil
}

func (openXmlService *ServiceImpl) GetUploadedOpenXmls(svc *cloudwatchlogs.CloudWatchLogs, queryParameters map[string]string) (int, *cloudwatchlogs.GetQueryResultsOutput) {
	if err := openXmlService.PrepareInsightsQueryParameters(queryParameters); err != nil {
		fmt.Println(err.Error())
		return http.StatusInternalServerError, nil
	}

	openXmlService.LogGroupName = "/aws/lambda/AWSUpload"

	queryExecutor := new(cloudwatch.QueryExecutorImpl)
	queryExecutor.Init(openXmlService.StartTimeEpoch, openXmlService.EndTimeEpoch, openXmlService.LogGroupName, openXmlService.QueryString)

	queryResultsOutput, err := queryExecutor.ExecuteQuery(svc)
	if err != nil {
		fmt.Println(err.Error())
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, queryResultsOutput
}
