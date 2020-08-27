package api

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"net/http"

	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudJson"
	myErrors "bitbucket.org/aldoft/printer-timeline-backend/app/internal/errors"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/heartbeat"
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

// ExtractQueryParams is responsible of extracting the query parameters from the gin context
// and returns a map with those query parameters.
func ExtractQueryParamsNew(r *events.APIGatewayProxyRequest) map[string]string {
	return map[string]string{
		"time_type":    r.QueryStringParameters["time_type"],
		"offset_units": r.QueryStringParameters["offset_units"],
		"offset_value": r.QueryStringParameters["offset_value"],
		"start_time":   r.QueryStringParameters["start_time"],
		"end_time":     r.QueryStringParameters["end_time"],
		"pn":           r.QueryStringParameters["pn"],
		"sn":           r.QueryStringParameters["sn"],
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

// ExtractStorageQueryParams is responsible of extracting the query parameters from the gin context
// and returns a map with those query parameters. It is used in the GetStoredObject endpoint.
func ExtractStorageQueryParamsNew(r *events.APIGatewayProxyRequest) map[string]string {
	return map[string]string{
		"bucket_region": r.QueryStringParameters["bucket_region"],
		"bucket_name":   r.QueryStringParameters["bucket_name"],
		"object_key":    r.QueryStringParameters["object_key"],
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
	xmlsFetcher openXml.OpenXmlsFetcher, cloudJsonsFetcher cloudJson.CloudJsonsFetcher,
	heartbeatsFetcher heartbeat.HeartbeatsFetcher) *gin.Engine {

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("api/object", StorageHandler(s3FetcherUsEast1, s3FetcherUsWest1))
	router.GET("api/open_xml", OpenXMLHandler(xmlsFetcher))
	router.GET("api/cloud_json", CloudJsonsHandler(cloudJsonsFetcher))
	router.GET("api/heartbeat", HeartbeatsHandler(heartbeatsFetcher))

	return router
}

func CreateLambdaHandler(s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher,
	xmlsFetcher DataFetcher, cloudJsonsFetcher DataFetcher,
	heartbeatsFetcher DataFetcher) LambdaHandler {

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
			handler = StorageHandlerNew(s3FetcherUsEast1, s3FetcherUsWest1)
		default:
			return newLambdaError(http.StatusBadRequest, myErrors.NotValidEndpoint)
		}

		return handler(ctx, request)
	}
}

func newLambdaError(httpStatus int, err error) (*events.APIGatewayProxyResponse, error) {

	return &events.APIGatewayProxyResponse{
		StatusCode: httpStatus,
		Body:       fmt.Sprintf("%s", err.Error()),
	}, nil

}

func newLambdaOkResponse(response []byte) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("%s", response),
	}, nil

}
