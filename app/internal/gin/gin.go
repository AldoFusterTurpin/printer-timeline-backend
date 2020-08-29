package gin

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/api"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/datafetcher"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/s3storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ExtractGinQueryParams is responsible of extracting the query parameters from the gin context
// and returns a map with those query parameters.
func ExtractGinQueryParams(c *gin.Context) map[string]string {
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

// ExtractGinStorageQueryParams is responsible of extracting the query parameters from the gin context
// and returns a map with those query parameters. It is used in the GetStoredObject endpoint.
func ExtractGinStorageQueryParams(c *gin.Context) map[string]string {
	return map[string]string{
		"bucket_region": c.Query("bucket_region"),
		"bucket_name":   c.Query("bucket_name"),
		"object_key":    c.Query("object_key"),
	}
}

// Handler is the responsible to handle the request.
// It returns a gin handler function that handles all the logic behind the http request.
// It uses dataFetcher interface that is responsible of fetching the data.
// Datafetcher, using polymorphism, will fetch the data depending on the implementation of dataFetcher
func Handler(dataFetcher datafetcher.DataFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryparams := ExtractGinQueryParams(c)
		status, result, err := api.GetData(queryparams, dataFetcher)

		if err != nil {
			c.JSON(status, err.Error())
		}
		c.JSON(status, result)
	}
}

// StorageHandler is the responsible to handle the request of get a specific object.
// It returns a gin handler function that handles all the logic behind the http request.
// It uses an s3fetcher interface that is responsible of fetching the stored objects (Openxml, CloudJson, HB, RTA, etc.).
// It calls GetStoredObject that is responsible of obtaiing the objects.
func StorageHandler(s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParams := ExtractGinStorageQueryParams(c)

		status, result, err := api.GetStoredObject(queryParams, s3FetcherUsEast1, s3FetcherUsWest1)

		if err != nil {
			c.JSON(status, err.Error())
		}
		c.JSON(status, result)
	}
}

// InitRouter initialize a gin router with all the routes for the different endpoints, request types and functions
// that are responsible of handling each request to specific endpoints.
func InitRouter(s3FetcherUsEast1 s3storage.S3Fetcher, s3FetcherUsWest1 s3storage.S3Fetcher,
	xmlsFetcher datafetcher.DataFetcher, cloudJsonsFetcher datafetcher.DataFetcher,
	heartbeatsFetcher datafetcher.DataFetcher, rtaFetcher datafetcher.RtaFetcher) *gin.Engine {

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("api/object", StorageHandler(s3FetcherUsEast1, s3FetcherUsWest1))
	router.GET("api/open_xml", Handler(xmlsFetcher))
	router.GET("api/cloud_json", Handler(cloudJsonsFetcher))
	router.GET("api/heartbeat", Handler(heartbeatsFetcher))
	router.GET("api/rta", Handler(rtaFetcher))

	return router
}
