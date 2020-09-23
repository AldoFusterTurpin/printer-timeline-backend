package queryparams

const (
	UsEast1S3Region = "US_EAST_1"
	UsWest1S3Region = "US_WEST_1"
)

// ExtractS3Info  extracts the AWS S3 information from the query parameters and
// returns the appropiate data and an error, if any.
func ExtractS3Info(queryParameters map[string]string) (bucketRegion string, bucketName string, objectKey string, err error) {
	bucketRegion = queryParameters["bucket_region"]
	bucketName = queryParameters["bucket_name"]
	objectKey = queryParameters["object_key"]

	if bucketRegion == "" {
		err = ErrorQueryStringMissingBucketRegion
		return
	}

	if bucketRegion != UsEast1S3Region && bucketRegion != UsWest1S3Region {
		err = ErrorQueryStringUnsupportedBucketRegion
		return
	}

	if bucketName == "" {
		err = ErrorQueryStringMissingBucketName
		return
	}

	if objectKey == "" {
		err = ErrorQueryStringMissingObjectKey
		return
	}

	return bucketRegion, bucketName, objectKey, nil
}
