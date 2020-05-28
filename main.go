package main

import (
	"bitbucket.org/aldoft/printer-timeline-backend/openXml"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-gonic/gin"
	"os"
)

func extractQueryParams(c *gin.Context) map[string]string{
	return map[string]string {
		"time_type" : c.Query("time_type"),
		"offset_units" : c.Query("offset_units"),
		"offset_value" : c.Query("offset_value"),
		"start_time" : c.Query("start_time"),
		"end_time" : c.Query("end_time"),
		"pn" : c.Query("pn"),
		"sn" : c.Query("sn"),
	}
}

func OpenXmlHandler(svc *cloudwatchlogs.CloudWatchLogs) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParameters := extractQueryParams(c)

		status, response := openXml.GetUploadedOpenXmls(svc, queryParameters)
		c.JSON(status, response)
	}
}

func initRouter(svc *cloudwatchlogs.CloudWatchLogs) *gin.Engine {
	router := gin.Default()
	router.GET("api/open_xml", OpenXmlHandler(svc))
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