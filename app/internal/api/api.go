package api

import (
	"net/http"

	myErrors "bitbucket.org/aldoft/printer-timeline-backend/app/internal/errors"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/openXml"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/s3storage"
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
		"bucket_name": c.Query("bucket_name"),
		"object_key":  c.Query("object_key"),
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
func InitRouter(s3fetcher s3storage.S3Fetcher, xmlsFetcher openXml.OpenXmlsFetcher) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("api/object", StorageHandler(s3fetcher))
	router.GET("api/open_xml", OpenXMLHandler(xmlsFetcher))

	return router
}
