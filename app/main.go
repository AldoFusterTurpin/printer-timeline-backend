package main

import (
	"errors"
	"fmt"
	"os"

	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/api"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudwatch"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/openXml"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/s3storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/s3"
)

func createAWSSession() (*session.Session, error) {
	envVarName := "AWS_REGION"

	awsRegion, ok := os.LookupEnv(envVarName)
	if !ok {
		return nil, errors.New("could not load " + envVarName + " environment variable")
	}

	return session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
}

func createCloudWatch(sess *session.Session) (*cloudwatchlogs.CloudWatchLogs, error) {
	sess, err := createAWSSession()
	if err != nil {
		return nil, err
	}
	return cloudwatchlogs.New(sess), nil
}

func createQueryExecutor(svc *cloudwatchlogs.CloudWatchLogs) cloudwatch.QueryExecutor {
	return cloudwatch.NewQueryExecutorImpl(svc)
}

//TODO: call this function in main and the result object will be used to handle the retrieval of
//the S3 data
func createS3Fetcher(sess *session.Session) s3storage.S3Fetcher {
	svc := s3.New(sess)
	return svc
}

func main() {
	sess, err := createAWSSession()
	if err != nil {
		fmt.Println(err)
		return
	}

	svc, err := createCloudWatch(sess)
	if err != nil {
		fmt.Println(err)
		return
	}

	queryExecutor := createQueryExecutor(svc)

	xmlsFetcher := openXml.NewOpenXmlsFetcherImpl(queryExecutor)

	s3Fetcher := createS3Fetcher(sess)

	router := api.InitRouter(s3Fetcher, xmlsFetcher)

	if err := router.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
