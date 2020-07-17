package api

import (
	"net/http"

	myErrors "bitbucket.org/aldoft/printer-timeline-backend/errors"
	"bitbucket.org/aldoft/printer-timeline-backend/openXml"
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

// SelectHttpStatus returns the appropiate http status based on the error passed as a parameter.
func SelectHttpStatus(err error) int {
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

// OpenXmlHandler is the responsible of handle the request of get the uploaded openXmls.
// It returns a gin handler function that handles all the logic behind the http request.
// It uses an xmlsFetcher interface that is responsible of fetching the OpenXMls.
func OpenXmlHandler(xmlsFetcher openXml.OpenXmlsFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParameters := ExtractQueryParams(c)

		result, err := xmlsFetcher.GetUploadedOpenXmls(queryParameters)

		status := SelectHttpStatus(err)

		if err != nil {
			c.JSON(status, err.Error())
		}
		c.JSON(status, result)
	}
}

// InitRouter initialize a gin router with all the routes for the different endpoints, request types and functions
// that are responsible of handling each request to specific endpoints.
func InitRouter(xmlsFetcher openXml.OpenXmlsFetcher) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("api/open_xml", OpenXmlHandler(xmlsFetcher))

	return router
}
