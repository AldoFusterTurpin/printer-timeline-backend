package awslambda

import (
	"context"
	"encoding/json"
	"net/http"

	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/api"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/datafetcher"
	myErrors "bitbucket.org/aldoft/printer-timeline-backend/app/internal/errors"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/s3storage"
	"github.com/aws/aws-lambda-go/events"
)

type LambdaHandler func(ctx context.Context, request *events.APIGatewayProxyRequest) (response *events.APIGatewayProxyResponse, err error)

func CreateLambdaHandler(s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher,
	xmlsFetcher datafetcher.DataFetcher, cloudJsonsFetcher datafetcher.DataFetcher,
	heartbeatsFetcher datafetcher.DataFetcher) LambdaHandler {

	return func(ctx context.Context, request *events.APIGatewayProxyRequest) (response *events.APIGatewayProxyResponse, err error) {
		var handler LambdaHandler

		switch request.Path {
		case "api/cloud_json":
			handler = GenericHandler(cloudJsonsFetcher)
		case "api/open_xml":
			handler = GenericHandler(xmlsFetcher)
		case "api/heartbeat":
			handler = GenericHandler(heartbeatsFetcher)
		case "api/object":
			handler = StorageHandler(s3FetcherUsEast1, s3FetcherUsWest1)
		default:
			return newLambdaError(http.StatusBadRequest, myErrors.NotValidEndpoint)
		}

		return handler(ctx, request)
	}
}

func GenericHandler(fetcher datafetcher.DataFetcher) LambdaHandler {
	return func(ctx context.Context, request *events.APIGatewayProxyRequest) (response *events.APIGatewayProxyResponse, err error) {
		queryParams := ExtractQueryParams(request)
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

func StorageHandler(s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher) LambdaHandler {
	return func(ctx context.Context, request *events.APIGatewayProxyRequest) (response *events.APIGatewayProxyResponse, err error) {
		queryParams := ExtractStorageQueryParams(request)

		status, result, err := api.GetStoredObject(queryParams, s3FetcherUsEast1, s3FetcherUsWest1)

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
