package queryparams

import "bitbucket.org/aldoft/printer-timeline-backend/app/internal/errors"

// ExtractS3Info  extracts the AWS S3 information from the query parameters and
// returns the appropiate data and an error, if any.
func ExtractS3Info(queryParameters map[string]string) (bucketRegion string, bucketName string, objectKey string, err error) {
	bucketRegion = queryParameters["bucket_region"]
	bucketName = queryParameters["bucket_name"]
	objectKey = queryParameters["object_key"]

	if bucketRegion == "" {
		err = errors.QueryStringMissingBucketRegion
		return
	}

	if bucketRegion != "US_EAST_1" && bucketRegion != "US_WEST_1" {
		err = errors.QueryStringUnsupportedBucketRegion
		return
	}

	if bucketName == "" {
		err = errors.QueryStringMissingBucketName
		return
	}

	if objectKey == "" {
		err = errors.QueryStringMissingObjectKey
		return
	}

	return bucketRegion, bucketName, objectKey, nil
}
