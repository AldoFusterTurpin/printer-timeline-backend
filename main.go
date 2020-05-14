package main

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func getListOfOpenXml(svc *cloudwatchlogs.CloudWatchLogs) gin.HandlerFunc {
	return func(c *gin.Context) {
		initialTime := c.DefaultQuery("from", "hello")
		finalTime := c.DefaultQuery("to", "world")

		c.String(http.StatusOK, "From: %s \nTo: %s", initialTime, finalTime)
	}
}

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := cloudwatchlogs.New(sess)


	router := gin.Default()
	router.GET("api/open_xml", getListOfOpenXml(svc))

	router.Run()
}

