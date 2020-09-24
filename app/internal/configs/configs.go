// Package configs provides the initialization of some parameters
package configs

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	EnvMaxTimeDiffMinutes   = "MAX_TIME_DIFF_IN_MINUTES"
	Development             = "DEVELOPMENT"
	AwsMainRegion           = "AWS_MAIN_REGION"
	AwsBlackseaBucketRegion = "AWS_BLACKSEA_BUCKET_REGION"
	DefaultTimeDiffMinutes  = 60
	MaxTimeDiffMinutes      = 2880

	InfraStructurePath = "/cc/V01/api/"
	CloudJsonPath      = InfraStructurePath + "cloud-json"
	OpenXMLPath        = InfraStructurePath + "open-xml"
	HeartbeatPath      = InfraStructurePath + "heartbeat"
	RTAPath            = InfraStructurePath + "rta"
	StorageObjectPath  = InfraStructurePath + "object"
	SubscriptionsPath  = InfraStructurePath + "subscriptions"

	ProductNumberQueryParam = "pn"
	SerialNumberQueryParam  = "sn"
	TimeTypeQueryParam      = "time_type"
	OffsetUnitsQueryParam   = "offset_units"
	OffsetValueQueryParam   = "offset_value"
	StartTimeQueryParam     = "start_time"
	EndTimeQueryParam       = "end_time"

	BucketRegionQueryParam = "bucket_region"
	BucketNameQueryParam   = "bucket_name"
	ObjectKeyQueryParam    = "object_key"
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

	stringDiff, ok := os.LookupEnv(EnvMaxTimeDiffMinutes)
	if !ok {
		return DefaultTimeDiffMinutes
	}

	intDiff, err := strconv.Atoi(stringDiff)
	if err != nil {
		return DefaultTimeDiffMinutes
	}

	if intDiff > MaxTimeDiffMinutes {
		return MaxTimeDiffMinutes
	}

	return intDiff
}

// Init initialises some global variables of configuration.
func Init() {
	maxTimeDiffInMinutes = setMaxTimeDiffInMinutes()
}

// IsDevelopment returns true if we are in development mode based in environment variables.
// Otherwise returns false.
func IsDevelopment() bool {

	dev, ok := os.LookupEnv(Development)
	if ok && strings.EqualFold(dev, "true") {
		return true
	}
	return false
}

// CreateAWSSession creates the corresponding sessions based on environment variables.
// It also returns an error, if any.
func CreateAWSSession() (sess1 *session.Session, sess2 *session.Session, err error) {

	awsRegion1, ok := os.LookupEnv(AwsMainRegion)
	if !ok {
		return nil, nil, errors.New("could not load " + AwsMainRegion + " environment variable")
	}

	sess1, err = session.NewSession(&aws.Config{
		Region: aws.String(awsRegion1)},
	)
	if err != nil {
		return nil, nil, err
	}

	awsRegion2, ok := os.LookupEnv(AwsBlackseaBucketRegion)
	if !ok {
		return nil, nil, errors.New("could not load " + AwsMainRegion + " environment variable")
	}

	sess2, err = session.NewSession(&aws.Config{
		Region: aws.String(awsRegion2)},
	)
	if err != nil {
		return nil, nil, err
	}

	return sess1, sess2, nil
}
