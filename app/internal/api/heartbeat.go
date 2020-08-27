package api

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/datafetcher"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-gonic/gin"
)

// GetHeartbeats is the responsible of obtaining the Heartbeats based in the queryParameters.
// This function is independent of the Framework used to create the web server as its input is just
// a map containing the http query parameters.
// A HeartbeatsFetcher is injected in order to obtain the Heartbeats.
func GetHeartbeats(queryParameters map[string]string, heartbeatsFetcher datafetcher.HeartbeatsFetcher) (status int, result *cloudwatchlogs.GetQueryResultsOutput, err error) {
	result, err = heartbeatsFetcher.GetUploadedHeartbeats(queryParameters)
	status = SelectHTTPStatus(err)
	return status, result, err
}

// HeartbeatsHandler is the responsible to handle the request of get the cloud Heartbeats.
// It returns a gin handler function that handles all the logic behind the http request.
// It uses an heartbeatsFetcher interface that is responsible of fetching the Heartbeats.
// It calls GetHeartbeats that is responsible of obtaiing the Xmls
func HeartbeatsHandler(heartbeatsFetcher datafetcher.HeartbeatsFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryparams := ExtractQueryParams(c)
		status, result, err := GetHeartbeats(queryparams, heartbeatsFetcher)

		if err != nil {
			c.JSON(status, err.Error())
		}
		c.JSON(status, result)
	}
}
