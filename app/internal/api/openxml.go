package api

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/openXml"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-gonic/gin"
)

// GetOpenXmls is the responsible of obtaining the Xmls based in the queryParameters.
// This function is independent of the Framework used to create the web server as its input is just
// a map containing the http query parameters.
// A xmlsFetcher is injected in order to obtain the Xmls.
func GetOpenXmls(queryParameters map[string]string, xmlsFetcher openXml.OpenXmlsFetcher) (status int, result *cloudwatchlogs.GetQueryResultsOutput, err error) {
	result, err = xmlsFetcher.GetUploadedOpenXmls(queryParameters)
	status = SelectHTTPStatus(err)
	return status, result, err

}

// OpenXMLHandler is the responsible to handle the request of get the uploaded openXmls.
// It returns a gin handler function that handles all the logic behind the http request.
// It uses an xmlsFetcher interface that is responsible of fetching the OpenXMls.
// It calls GetOpenXmls that is responsible of obtaiing the Xmls
func OpenXMLHandler(xmlsFetcher openXml.OpenXmlsFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryparams := ExtractQueryParams(c)
		status, result, err := GetOpenXmls(queryparams, xmlsFetcher)

		if err != nil {
			c.JSON(status, err.Error())
		}
		c.JSON(status, result)
	}
}
