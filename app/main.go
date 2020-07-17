package main

import (
	"errors"
	"fmt"
	"os"

	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/api"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudwatch"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/openXml"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func createCloudWatchServiceClient() (*cloudwatchlogs.CloudWatchLogs, error) {
	envVarName := "AWS_REGION"

	awsRegion, ok := os.LookupEnv(envVarName)
	if !ok {
		return nil, errors.New("could not load " + envVarName + " environment variable")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	if err != nil {
		return nil, err
	}

	return cloudwatchlogs.New(sess), nil
}

func createQueryExecutor() (cloudwatch.QueryExecutor, error) {
	svc, err := createCloudWatchServiceClient()
	if err != nil {
		return nil, err
	}

	return cloudwatch.NewQueryExecutorImpl(svc), nil
}

func main() {
	queryExecutor, err := createQueryExecutor()
	if err != nil {
		fmt.Println(err)
		return
	}

	xmlsFetcher := openXml.NewOpenXmlsFetcherImpl(queryExecutor)

	router := api.InitRouter(xmlsFetcher)

	if err := router.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
