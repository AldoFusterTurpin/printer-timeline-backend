// Package init provides the initialization of some parameters
package init

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	maxTimeDiffInMinutes int
)

// GetMaxTimeDiffInMinutes returns the maximum amount of time difference (in minutes)
// between start_time and end_time allowed in a query.
func GetMaxTimeDiffInMinutes() int {
	return maxTimeDiffInMinutes
}

func setMaxTimeDiffInMinutes() int {
	stringDiff, ok := os.LookupEnv("MAX_TIME_DIFF_IN_MINUTES")
	if !ok {
		return 60
	}

	intDiff, err := strconv.Atoi(stringDiff)
	if err != nil {
		return 60
	}

	// the maximum time difference allowed for queries.
	twoWeeksInMinutes := 20160
	if intDiff > twoWeeksInMinutes {
		return twoWeeksInMinutes
	}
	return intDiff
}

// Init initialises some global variables of configuration.
func Init() {
	maxTimeDiffInMinutes = setMaxTimeDiffInMinutes()
}

// IsDevelopment returns true if we are in deveolpment mode based in environment variables.
// Otherwise returns false.
func IsDevelopment() bool {
	dev, ok := os.LookupEnv("DEVELOPMENT")
	if ok && strings.EqualFold(dev, "true") {
		return true
	}
	return false
}

// CreateAWSSession creates the corresponding sessions based on environment variables.
// It also returns an error, if any.
func CreateAWSSession() (sess1 *session.Session, sess2 *session.Session, err error) {
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
