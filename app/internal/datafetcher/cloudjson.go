package datafetcher

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// CloudJsonsFetcher is the implementation of datasFetcher that uses a queryExecutor to perform a query
//and obtain the Cloud Jsons.
type CloudJsonsFetcher struct {
	queryExecutor cloudwatch.QueryExecutor
}

// NewCloudJsonsFetcher creates a new CloudJsonsFetcher implementationm
func NewCloudJsonsFetcher(queryExecutor cloudwatch.QueryExecutor) CloudJsonsFetcher {
	return CloudJsonsFetcher{queryExecutor}
}

// GetLogGroupName returns the appropiate Log group in AWS CloudWatch
func (cloudJsonsFetcher CloudJsonsFetcher) GetLogGroupName() (logGroupName string) {
	return "/aws/lambda/AWSParser"
}

// CreateQueryTemplate returns a new query template depending on the productNumber and serialNumber parameters.
// The resulting query template will be used by the queryExecutor to obtain the Cloud Jsons.
func (cloudJsonsFetcher CloudJsonsFetcher) CreateQueryTemplate(productNumber, serialNumber string) (queryTemplateString string) {
	// As the query string contains backticks ``, I need to surround those backticks with double quotes("").
	// There is also the expression "{{.productNumber}}" which need to be
	// surrounded by double quotes ("") (because the Go backticks used to create a string with multiple lines can't contain double quotes inside it).
	// Last, I need to concatenate all those strings but need to create temporary variables because are 'untyped string constants'.
	// For more info, check: https://blog.golang.org/constants
	if productNumber != "" && serialNumber != "" {
		s1 := `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date, `
		s2 := "`fields.metadata.xml-generator-object-path`"
		s3 := `| filter (ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date)`
		s4 := " and ispresent(`fields.metadata.xml-generator-object-path`)"
		s5 := ` and fields.topic = "json" and fields.ProductNumber="{{.productNumber}}" and fields.SerialNumber="{{.serialNumber}}")
		| sort @timestamp asc
		| limit 10000`

		return s1 + s2 + s3 + s4 + s5
	}
	if productNumber != "" {
		s1 := `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date, `
		s2 := "`fields.metadata.xml-generator-object-path`"
		s3 := `| filter (ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date)`
		s4 := " and ispresent(`fields.metadata.xml-generator-object-path`)"
		s5 := ` and fields.topic = "json" and fields.ProductNumber="{{.productNumber}}")
		| sort @timestamp asc
		| limit 10000`

		return s1 + s2 + s3 + s4 + s5
	}

	s1 := `fields @timestamp, fields.ProductNumber, fields.SerialNumber, fields.bucket_name, fields.bucket_region, fields.key, fields.topic, fields.metadata.date, `
	s2 := "`fields.metadata.xml-generator-object-path`"
	s3 := `| filter (ispresent(fields.ProductNumber) and ispresent(fields.SerialNumber) and ispresent(fields.bucket_name) and ispresent(fields.bucket_region) and ispresent(fields.key) and ispresent(fields.topic) and ispresent(fields.metadata.date)`
	s4 := " and ispresent(`fields.metadata.xml-generator-object-path`)"
	s5 := ` and fields.topic = "json")`
	s6 := `| sort @timestamp asc | limit 10000`

	return s1 + s2 + s3 + s4 + s5 + s6
}

// FetchData obtains the Jsons created by Cloud Connector depending on requestQueryParams.
// The method basically creates a variable insightsQueryParams and then calls a queryExecutor
// to perform the query. It returns the result and an error, if any.
func (cloudJsonsFetcher CloudJsonsFetcher) FetchData(requestQueryParams map[string]string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	insightsQueryParams, err := createInsightsQueryParams(requestQueryParams, cloudJsonsFetcher)
	if err != nil {
		return nil, err
	}

	result, err := cloudJsonsFetcher.queryExecutor.ExecuteQuery(insightsQueryParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}
