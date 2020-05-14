package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func convertEpochStringToUint64(epochToConvert string, defaultEpoch int64) (epochConverted int64, err error) {
	if epochToConvert == "" {
		return defaultEpoch, nil
	}
	return strconv.ParseInt(epochToConvert, 10, 64)
}


func defaultStartTime() time.Time {
	return time.Now().Add(-time.Minute * 15)
}


func defaultEndTime() time.Time {
	return time.Now()
}


func queryUplodadedOpenXmls(svc *cloudwatchlogs.CloudWatchLogs) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var startTimeEpoch, endTimeEpoch int64

		startTime := c.Query("startTime")
		startTimeEpoch, err = convertEpochStringToUint64(startTime, defaultStartTime().Unix())
		 if err != nil {
			fmt.Println(err.Error())
			return
		}

		endTime := c.Query("endTime")
		endTimeEpoch, err = convertEpochStringToUint64(endTime, defaultEndTime().Unix())
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		queryString := `fields @timestamp, @message
		| sort @timestamp desc
		| limit 20`

		logGroupName := "/aws/lambda/AWSUpload"

		input := &cloudwatchlogs.StartQueryInput{
			EndTime: &endTimeEpoch,
			LogGroupName: aws.String(logGroupName),
			QueryString: aws.String(queryString),
			StartTime: &startTimeEpoch,
		}
		var startQueryOutput *cloudwatchlogs.StartQueryOutput
		startQueryOutput, err = svc.StartQuery(input)



		c.String(http.StatusOK, "Hi")
	}
}

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := cloudwatchlogs.New(sess)

	router := gin.Default()
	router.GET("api/open_xml", queryUplodadedOpenXmls(svc))

	err := router.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}

