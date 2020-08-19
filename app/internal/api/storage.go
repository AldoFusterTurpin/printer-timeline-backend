package api

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/s3storage"
	"github.com/gin-gonic/gin"
)

//TODO: add region to back-end api

// GetStoredObject is the responsible of obtaining the stored objects based in the queryParameters.
// This function is independent of the Framework used to create the web server as its input is just
// a map containing the http query parameters.
// A xmlsFetcher is injected in order to obtain the stored objects.
func GetStoredObject(queryParameters map[string]string, s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher) (status int, result string, err error) {
	var bytesResult []byte

	if queryParameters["bucket_region"] == "US_EAST_1" {
		bytesResult, err = s3storage.GetS3Data(s3FetcherUsEast1, queryParameters["bucket_name"], queryParameters["object_key"])
	} else {
		bytesResult, err = s3storage.GetS3Data(s3FetcherUsWest1, queryParameters["bucket_name"], queryParameters["object_key"])
	}

	status = SelectHTTPStatus(err)
	return status, string(bytesResult), err
}

// StorageHandler is the responsible to handle the request of get a specific object.
// It returns a gin handler function that handles all the logic behind the http request.
// It uses an s3fetcher interface that is responsible of fetching the stored objects (Openxml, CloudJson, HB, RTA, etc.).
// It calls GetStoredObject that is responsible of obtaiing the objects.
func StorageHandler(s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParams := ExtractStorageQueryParams(c)

		status, result, err := GetStoredObject(queryParams, s3FetcherUsEast1, s3FetcherUsWest1)

		if err != nil {
			c.JSON(status, err.Error())
		}
		c.JSON(status, result)
	}
}
