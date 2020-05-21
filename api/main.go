package main

import (
	"bitbucket.org/aldoft/printer-timeline-backend/api/openXml"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-gonic/gin"
	"os"
)

func initRouter(svc *cloudwatchlogs.CloudWatchLogs) *gin.Engine {
	router := gin.Default()
	router.GET("api/open_xml", openXml.GetUploadedOpenXmls(svc))
	return router
}

func main() {
	awsRegion, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		fmt.Println(errors.New("could not load BUCKET_REGION env var"))
		return
	}

	var err error
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	svc := cloudwatchlogs.New(sess)

	router := initRouter(svc)

	if err = router.Run(); err != nil {
		fmt.Println(err.Error())
		return
	}
}