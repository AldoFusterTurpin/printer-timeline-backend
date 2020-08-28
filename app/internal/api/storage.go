package api

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/queryparams"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/s3storage"
)

// GetStoredObject is the responsible of obtaining the stored objects based in the queryParameters.
// This function is independent of the Framework used to create the web server as its input is just
// a map containing the http query parameters.
// A xmlsFetcher is injected in order to obtain the stored objects.
func GetStoredObject(queryParameters map[string]string, s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher) (status int, result string, err error) {
	bucketRegion, bucketName, objectKey, err := queryparams.ExtractS3Info(queryParameters)
	if err != nil {
		status = SelectHTTPStatus(err)
		return
	}

	s3FetcherToUSe := selectS3Fetcher(bucketRegion, s3FetcherUsEast1, s3FetcherUsWest1)

	var bytesResult []byte
	bytesResult, err = s3storage.GetS3Data(*s3FetcherToUSe, bucketName, objectKey)

	status = SelectHTTPStatus(err)
	return status, string(bytesResult), err
}

func selectS3Fetcher(bucketRegion string, s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher) *s3storage.S3Fetcher {
	if bucketRegion == "US_EAST_1" {
		return &s3FetcherUsEast1
	}
	if bucketRegion == "US_WEST_1" {
		return &s3FetcherUsWest1
	}
	return nil
}
