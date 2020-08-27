package api

import (
	"net/http"

	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/datafetcher"
	myErrors "bitbucket.org/aldoft/printer-timeline-backend/app/internal/errors"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/s3storage"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ExtractQueryParams is responsible of extracting the query parameters from the gin context
// and returns a map with those query parameters.
func ExtractQueryParams(c *gin.Context) map[string]string {
	return map[string]string{
		"time_type":    c.Query("time_type"),
		"offset_units": c.Query("offset_units"),
		"offset_value": c.Query("offset_value"),
		"start_time":   c.Query("start_time"),
		"end_time":     c.Query("end_time"),
		"pn":           c.Query("pn"),
		"sn":           c.Query("sn"),
	}
}

// ExtractStorageQueryParams is responsible of extracting the query parameters from the gin context
// and returns a map with those query parameters. It is used in the GetStoredObject endpoint.
func ExtractStorageQueryParams(c *gin.Context) map[string]string {
	return map[string]string{
		"bucket_region": c.Query("bucket_region"),
		"bucket_name":   c.Query("bucket_name"),
		"object_key":    c.Query("object_key"),
	}
}

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

// InitRouter initialize a gin router with all the routes for the different endpoints, request types and functions
// that are responsible of handling each request to specific endpoints.
func InitRouter(s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher,
	xmlsFetcher datafetcher.DataFetcher, cloudJsonsFetcher datafetcher.DataFetcher,
	heartbeatsFetcher datafetcher.DataFetcher) *gin.Engine {

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("api/object", StorageHandler(s3FetcherUsEast1, s3FetcherUsWest1))
	router.GET("api/open_xml", Handler(xmlsFetcher))
	router.GET("api/cloud_json", Handler(cloudJsonsFetcher))
	router.GET("api/heartbeat", Handler(heartbeatsFetcher))

	return router
}

// GetData is the responsible of obtaining the Jsons based in the queryParameters.
// This function is independent of the Framework used to create the web server as its input is just
// a map containing the http query parameters.
// A cloudJsonFetcher is injected in order to obtain the Jsons.
func GetData(queryParameters map[string]string, fetcher datafetcher.DataFetcher) (status int, result *cloudwatchlogs.GetQueryResultsOutput, err error) {
	result, err = fetcher.FetchData(queryParameters)
	status = SelectHTTPStatus(err)
	return status, result, err

}

// CloudJsonsHandler is the responsible to handle the request of get the cloud Jsons.
// It returns a gin handler function that handles all the logic behind the http request.
// It uses an cloudJsonsFetcher interface that is responsible of fetching the Jsons.
// It calls GetData because is responsible of obtaiing the Xmls
func Handler(dataFetcher datafetcher.DataFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryparams := ExtractQueryParams(c)
		status, result, err := GetData(queryparams, dataFetcher)

		if err != nil {
			c.JSON(status, err.Error())
		}
		c.JSON(status, result)
	}
}
