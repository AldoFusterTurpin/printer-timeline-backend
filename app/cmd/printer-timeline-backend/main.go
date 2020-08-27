package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"strings"

	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/heartbeat"

	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/api"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudJson"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudwatch"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/openXml"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/s3storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/s3"
)

func createAWSSession() (sess1 *session.Session, sess2 *session.Session, err error) {

	envVarName := "MAIN_AWS_REGION"
	envVarName2 := "AWS_REGION_BLACK_SEA_BUCKET"

	awsRegion1, ok := os.LookupEnv(envVarName)
	if !ok {
		return nil, nil, errors.New("could not load " + envVarName + " environment variable")
	}

	sess1, err = session.NewSession(&aws.Config{
		Region: aws.String(awsRegion1)},
	)
	if err != nil {
		return nil, nil, err
	}

	awsRegion2, ok := os.LookupEnv(envVarName2)
	if !ok {
		return nil, nil, errors.New("could not load " + envVarName + " environment variable")
	}

	sess2, err = session.NewSession(&aws.Config{
		Region: aws.String(awsRegion2)},
	)
	if err != nil {
		return nil, nil, err
	}

	return sess1, sess2, nil
}

func createCloudWatch(sess *session.Session) (*cloudwatchlogs.CloudWatchLogs, error) {
	return cloudwatchlogs.New(sess), nil
}

func createQueryExecutor(svc *cloudwatchlogs.CloudWatchLogs) cloudwatch.QueryExecutor {
	return cloudwatch.NewQueryExecutorImpl(svc)
}

func createS3Fetcher(sess *session.Session) s3storage.S3Fetcher {
	svc := s3.New(sess)
	return svc
}

func main() {
	sess1, sess2, err := createAWSSession()
	if err != nil {
		fmt.Println(err)
		return
	}

	svc, err := createCloudWatch(sess1)
	if err != nil {
		fmt.Println(err)
		return
	}

	queryExecutor := createQueryExecutor(svc)

	xmlsFetcher := openXml.NewOpenXmlsFetcherImpl(queryExecutor)
	cloudJsonFetcher := cloudJson.NewCloudJsonsFetcherImpl(queryExecutor)
	heartbeatsFetcher := heartbeat.NewHeartbeatsFetcherImpl(queryExecutor)

	s3FetcherUsEast1 := createS3Fetcher(sess1)
	s3FetcherUsWest1 := createS3Fetcher(sess2)

	dev := isDevelopment()
	if dev {
		router := api.InitRouter(s3FetcherUsEast1, s3FetcherUsWest1, xmlsFetcher, cloudJsonFetcher, heartbeatsFetcher)
		if err := router.Run(); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		lambda.Start(api.CreateLambdaHandler(s3FetcherUsEast1, s3FetcherUsWest1, xmlsFetcher, cloudJsonFetcher, heartbeatsFetcher))
	}
}

func isDevelopment() bool {
	dev, ok := os.LookupEnv("DEVELOPMENT")
	if ok && strings.EqualFold(dev, "true") {
		return true
	}

	return false
}
