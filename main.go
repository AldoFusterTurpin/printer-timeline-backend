package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"bitbucket.org/aldoft/printer-timeline-backend/openXml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func extractQueryParams(c *gin.Context) map[string]string {
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

func selectHttpStatus(err error) int {
	switch err.(type) {
	case nil:
		return http.StatusOK
	/* case yourError:

	case otherError:
	*/
	default:
		return http.StatusInternalServerError
	}
}

func OpenXmlHandler(xmlsFetcher openXml.OpenXmlsFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParameters := extractQueryParams(c)

		result, err := xmlsFetcher.GetUploadedOpenXmls(queryParameters)

		status := selectHttpStatus(err)

		if err != nil {
			c.JSON(status, result)
		} else {
			c.JSON(status, result)
		}
	}
}

func initRouter(xmlsFetcher openXml.OpenXmlsFetcher) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("api/open_xml", OpenXmlHandler(xmlsFetcher))

	return router
}

func newCloudWatchService() (*cloudwatchlogs.CloudWatchLogs, error) {
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

	svc := cloudwatchlogs.New(sess)
	return svc, nil
}

func main() {
	svc, err := newCloudWatchService()
	if err != nil {
		fmt.Println(err)
		return
	}

	xmlsFetcher := openXml.NewOpenXmlsFetcherImpl(svc)
	router := initRouter(xmlsFetcher)

	if err := router.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
