package api

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/queryparams"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/s3storage"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetStoredObject is the responsible of obtaining the stored objects based in the queryParameters.
// This function is independent of the Framework used to create the web server as its input is just
// a map containing the http query parameters.
// A xmlsFetcher is injected in order to obtain the stored objects.
func GetStoredObject(queryParameters map[string]string, s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher) (status int, result string, err error) {
	bucketRegion, bucketName, objectKey, err := queryparams.ExtractS3Info(queryParameters)
	if err != nil {
		status = SelectHTTPStatus(err)
		return
	}

	var bytesResult []byte

	s3FetcherToUSe := selectS3Fetcher(bucketRegion, s3FetcherUsEast1, s3FetcherUsWest1)

	bytesResult, err = s3storage.GetS3Data(*s3FetcherToUSe, bucketName, objectKey)

	status = SelectHTTPStatus(err)
	return status, string(bytesResult), err
}

func selectS3Fetcher(bucketRegion string, s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher) *s3storage.S3Fetcher {
	if bucketRegion == "US_EAST_1" {
		return &s3FetcherUsEast1
	}
	if bucketRegion == "US_WEST_1" {
		return &s3FetcherUsWest1
	}
	return nil
}

// StorageHandler is the responsible to handle the request of get a specific object.
// It returns a gin handler function that handles all the logic behind the http request.
// It uses an s3fetcher interface that is responsible of fetching the stored objects (Openxml, CloudJson, HB, RTA, etc.).
// It calls GetStoredObject that is responsible of obtaiing the objects.
func StorageHandler(s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParams := ExtractStorageQueryParams(c)

		status, result, err := GetStoredObject(queryParams, s3FetcherUsEast1, s3FetcherUsWest1)

		if err != nil {
			c.JSON(status, err.Error())
		}
		c.JSON(status, result)
	}
}

func StorageHandlerNew(s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher) LambdaHandler {
	return func(ctx context.Context, request *events.APIGatewayProxyRequest) (response *events.APIGatewayProxyResponse, err error) {
		queryParams := ExtractStorageQueryParamsNew(request)

		status, result, err := GetStoredObject(queryParams, s3FetcherUsEast1, s3FetcherUsWest1)

		if err != nil {
			return newLambdaError(status, err)
		}

		jsonResp, err := json.Marshal(result)
		if err != nil {
			return newLambdaError(http.StatusInternalServerError, err)
		}

		return newLambdaOkResponse(jsonResp)
	}
}
