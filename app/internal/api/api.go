package api

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/db"
	"net/http"

	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/datafetcher"
	. "bitbucket.org/aldoft/printer-timeline-backend/app/internal/queryparams"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

//go:generate mockgen -source=../datafetcher/datafetcher.go -destination=../datafetcher/mocks/datafetcher.go -package=mocks

// SelectHTTPStatus returns the appropriate http status based on the error passed as a parameter.
func SelectHTTPStatus(err error) int {
	switch err {
	case nil:
		return http.StatusOK
	case ErrorQueryStringMissingTimeRangeType, ErrorQueryStringUnsupportedTimeRangeType, ErrorQueryStringStartTimeAppears,
		ErrorQueryStringMissingEndTime, ErrorQueryStringEndTimeAppears, ErrorQueryStringUnsupportedEndTime,
		ErrorQueryStringMissingOffsetUnits, ErrorQueryStringUnsupportedOffsetUnits, ErrorQueryStringMissingOffsetValue,
		ErrorQueryStringUnsupportedOffsetValue, ErrorQueryStringMissingStartTime, ErrorQueryStringUnsupportedStartTime,
		ErrorQueryStringTimeDifferenceTooBig, ErrorQueryStringEndTimePreviousThanStartTime, ErrorQueryStringPnSn:
		return http.StatusBadRequest
	case db.NotFoundErr:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// GetData is the responsible of obtaining the data based in the queryParameters.
// This function is independent of the Framework used to create the web server as its input is just
// a maputil containing the http query parameters.
// A DataFetcher is injected in order to obtain the data.
func GetData(queryParameters map[string]string, fetcher datafetcher.DataFetcher) (status int, result *cloudwatchlogs.GetQueryResultsOutput, err error) {
	result, err = fetcher.FetchData(queryParameters)
	status = SelectHTTPStatus(err)
	return status, result, err
}
