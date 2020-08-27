package awslambda

import "github.com/aws/aws-lambda-go/events"

// ExtractQueryParams is responsible of extracting the query parameters from the
// API Gateway request.
func ExtractQueryParams(r *events.APIGatewayProxyRequest) map[string]string {
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

// ExtractStorageQueryParams is responsible of extracting the query parameters from the
// API Gateway request and returns a map with those query parameters.
func ExtractStorageQueryParams(r *events.APIGatewayProxyRequest) map[string]string {
	return map[string]string{
		"bucket_region": r.QueryStringParameters["bucket_region"],
		"bucket_name":   r.QueryStringParameters["bucket_name"],
		"object_key":    r.QueryStringParameters["object_key"],
	}
}
