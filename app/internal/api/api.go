package api

import (
	"net/http"

	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/datafetcher"
	myErrors "bitbucket.org/aldoft/printer-timeline-backend/app/internal/errors"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

//go:generate mockgen -source=../datafetcher/datafetcher.go -destination=../datafetcher/mocks/datafetcher.go -package=mocks

// SelectHTTPStatus returns the appropiate http status based on the error passed as a parameter.
func SelectHTTPStatus(err error) int {
	switch err {
	case nil:
		return http.StatusOK
	case myErrors.QueryStringMissingTimeRangeType, myErrors.QueryStringUnsupportedTimeRangeType, myErrors.QueryStringStartTimeAppears,
		myErrors.QueryStringMissingEndTime, myErrors.QueryStringEndTimeAppears, myErrors.QueryStringUnsupportedEndTime,
		myErrors.QueryStringMissingOffsetUnits, myErrors.QueryStringUnsupportedOffsetUnits, myErrors.QueryStringMissingOffsetValue,
		myErrors.QueryStringUnsupportedOffsetValue, myErrors.QueryStringMissingStartTime, myErrors.QueryStringUnsupportedStartTime,
		myErrors.QueryStringTimeDifferenceTooBig, myErrors.QueryStringEndTimePreviousThanStartTime, myErrors.QueryStringPnSn:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// GetData is the responsible of obtaining the data based in the queryParameters.
// This function is independent of the Framework used to create the web server as its input is just
// a map containing the http query parameters.
// A DataFetcher is injected in order to obtain the data.
func GetData(queryParameters map[string]string, fetcher datafetcher.DataFetcher) (status int, result *cloudwatchlogs.GetQueryResultsOutput, err error) {
	result, err = fetcher.FetchData(queryParameters)
	status = SelectHTTPStatus(err)
	return status, result, err
}
