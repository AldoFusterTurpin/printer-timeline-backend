package api

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/datafetcher"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-gonic/gin"
)

// GetCloudJsons is the responsible of obtaining the Jsons based in the queryParameters.
// This function is independent of the Framework used to create the web server as its input is just
// a map containing the http query parameters.
// A cloudJsonFetcher is injected in order to obtain the Jsons.
func GetCloudJsons(queryParameters map[string]string, fetcher datafetcher.DataFetcher) (status int, result *cloudwatchlogs.GetQueryResultsOutput, err error) {
	result, err = fetcher.FetchData(queryParameters)
	status = SelectHTTPStatus(err)
	return status, result, err

}

// CloudJsonsHandler is the responsible to handle the request of get the cloud Jsons.
// It returns a gin handler function that handles all the logic behind the http request.
// It uses an cloudJsonsFetcher interface that is responsible of fetching the Jsons.
// It calls GetCloudJsons that is responsible of obtaiing the Xmls
func CloudJsonsHandler(cloudJsonsFetcher datafetcher.CloudJsonsFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryparams := ExtractQueryParams(c)
		status, result, err := GetCloudJsons(queryparams, cloudJsonsFetcher)

		if err != nil {
			c.JSON(status, err.Error())
		}
		c.JSON(status, result)
	}
}
