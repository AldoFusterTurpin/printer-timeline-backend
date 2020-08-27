package api

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudJson"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetCloudJsons is the responsible of obtaining the Jsons based in the queryParameters.
// This function is independent of the Framework used to create the web server as its input is just
// a map containing the http query parameters.
// A cloudJsonFetcher is injected in order to obtain the Jsons.
func GetCloudJsons(queryParameters map[string]string, fetcher DataFetcher) (status int, result *cloudwatchlogs.GetQueryResultsOutput, err error) {
	result, err = fetcher.FetchData(queryParameters)
	status = SelectHTTPStatus(err)
	return status, result, err

}

// CloudJsonsHandler is the responsible to handle the request of get the cloud Jsons.
// It returns a gin handler function that handles all the logic behind the http request.
// It uses an cloudJsonsFetcher interface that is responsible of fetching the Jsons.
// It calls GetCloudJsons that is responsible of obtaiing the Xmls
func CloudJsonsHandler(cloudJsonsFetcher cloudJson.CloudJsonsFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryparams := ExtractQueryParams(c)
		status, result, err := GetCloudJsons(queryparams, cloudJsonsFetcher)

		if err != nil {
			c.JSON(status, err.Error())
		}
		c.JSON(status, result)
	}
}

type LambdaHandler func(ctx context.Context, request *events.APIGatewayProxyRequest) (response *events.APIGatewayProxyResponse, err error)

func GenericHandler(fetcher DataFetcher) LambdaHandler {
	return func(ctx context.Context, request *events.APIGatewayProxyRequest) (response *events.APIGatewayProxyResponse, err error) {
		queryParams := ExtractQueryParamsNew(request)
		result, err := fetcher.FetchData(queryParams)

		if err != nil {
			return newLambdaError(http.StatusInternalServerError, err)
		}

		jsonResp, err := json.Marshal(result)
		if err != nil {
			return newLambdaError(http.StatusInternalServerError, err)
		}

		return newLambdaOkResponse(jsonResp)
	}
}

type DataFetcher interface {
	FetchData(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error)
}
