package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"

	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/awslambda"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/cloudwatch"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/datafetcher"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/gin"
	initConfig "bitbucket.org/aldoft/printer-timeline-backend/app/internal/init"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/s3storage"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/s3"
)

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
	initConfig.Init()

	sess1, sess2, err := initConfig.CreateAWSSession()
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

	xmlsFetcher := datafetcher.NewOpenXmlsFetcher(queryExecutor)
	cloudJsonFetcher := datafetcher.NewCloudJsonsFetcher(queryExecutor)
	heartbeatsFetcher := datafetcher.NewHeartbeatsFetcher(queryExecutor)
	rtaFetcher := datafetcher.NewRtaFetcher(queryExecutor)

	s3FetcherUsEast1 := createS3Fetcher(sess1)
	s3FetcherUsWest1 := createS3Fetcher(sess2)

	dev := initConfig.IsDevelopment()
	if dev {
		router := gin.InitRouter(s3FetcherUsEast1, s3FetcherUsWest1, xmlsFetcher, cloudJsonFetcher, heartbeatsFetcher, rtaFetcher)
		if err := router.Run(); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		lambda.Start(awslambda.CreateLambdaHandler(s3FetcherUsEast1, s3FetcherUsWest1, xmlsFetcher, cloudJsonFetcher, heartbeatsFetcher, rtaFetcher))
	}
}
