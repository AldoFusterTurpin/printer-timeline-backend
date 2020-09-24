package awslambda

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/configs"
	"github.com/aws/aws-lambda-go/events"
)

// ExtractQueryParams is responsible of extracting the query parameters from the
// API Gateway request.
func ExtractQueryParams(r *events.APIGatewayProxyRequest) map[string]string {
	return map[string]string{
		configs.TimeTypeQueryParam:      r.QueryStringParameters[configs.TimeTypeQueryParam],
		configs.OffsetUnitsQueryParam:   r.QueryStringParameters[configs.OffsetUnitsQueryParam],
		configs.OffsetValueQueryParam:   r.QueryStringParameters[configs.OffsetValueQueryParam],
		configs.StartTimeQueryParam:     r.QueryStringParameters[configs.StartTimeQueryParam],
		configs.EndTimeQueryParam:       r.QueryStringParameters[configs.EndTimeQueryParam],
		configs.ProductNumberQueryParam: r.QueryStringParameters[configs.ProductNumberQueryParam],
		configs.SerialNumberQueryParam:  r.QueryStringParameters[configs.SerialNumberQueryParam],
	}
}

// ExtractStorageQueryParams is responsible of extracting the query parameters from the
// API Gateway request and returns a map with those query parameters.
func ExtractStorageQueryParams(r *events.APIGatewayProxyRequest) map[string]string {
	return map[string]string{
		configs.BucketRegionQueryParam: r.QueryStringParameters[configs.BucketRegionQueryParam],
		configs.BucketNameQueryParam:   r.QueryStringParameters[configs.BucketNameQueryParam],
		configs.ObjectKeyQueryParam:    r.QueryStringParameters[configs.ObjectKeyQueryParam],
	}
}

// ExtractPrinterQueryParams is responsible of extracting the printer query parameters from the
// API Gateway request and return a map with them
func ExtractPrinterQueryParams(r *events.APIGatewayProxyRequest) map[string]string {
	return map[string]string{
		configs.ProductNumberQueryParam: r.QueryStringParameters[configs.ProductNumberQueryParam],
		configs.SerialNumberQueryParam:  r.QueryStringParameters[configs.SerialNumberQueryParam],
	}
}
