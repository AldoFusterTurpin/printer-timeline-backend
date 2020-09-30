package gin

import (
	"net/http"
	"time"

	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/api"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/configs"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/datafetcher"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/db"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/maputil"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ExtractGinPrinterQueryParams is responsible of extracting the query parameters from the gin context
// and returns a maputil with those query parameters.
func ExtractGinPrinterQueryParams(c *gin.Context) map[string]string {
	return map[string]string{
		configs.ProductNumberQueryParam: c.Query(configs.ProductNumberQueryParam),
		configs.SerialNumberQueryParam:  c.Query(configs.SerialNumberQueryParam),
	}
}

// ExtractGinQueryParams is responsible of extracting the query parameters from the gin context
// and returns a maputil with those query parameters.
func extractSpecificQueryParams(c *gin.Context) map[string]string {

	return map[string]string{
		configs.TimeTypeQueryParam:    c.Query(configs.TimeTypeQueryParam),
		configs.OffsetUnitsQueryParam: c.Query(configs.OffsetUnitsQueryParam),
		configs.OffsetValueQueryParam: c.Query(configs.OffsetValueQueryParam),
		configs.StartTimeQueryParam:   c.Query(configs.StartTimeQueryParam),
		configs.EndTimeQueryParam:     c.Query(configs.EndTimeQueryParam),
	}
}

// ExtractGinQueryParams is responsible of extracting the query parameters from the gin context
// and returns a maputil with those query parameters.
func ExtractGinQueryParams(c *gin.Context) map[string]string {
	printerQueryParams := ExtractGinPrinterQueryParams(c)
	specificQueryParams := extractSpecificQueryParams(c)

	return maputil.JoinMaps(printerQueryParams, specificQueryParams)
}

// ExtractGinStorageQueryParams is responsible of extracting the query parameters from the gin context
// and returns a maputil with those query parameters. It is used in the GetStoredObject endpoint.
func ExtractGinStorageQueryParams(c *gin.Context) map[string]string {
	return map[string]string{
		configs.BucketRegionQueryParam: c.Query(configs.BucketRegionQueryParam),
		configs.BucketNameQueryParam:   c.Query(configs.BucketNameQueryParam),
		configs.ObjectKeyQueryParam:    c.Query(configs.ObjectKeyQueryParam),
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
		} else {
			c.JSON(status, result)
		}
	}
}

// StorageHandler is the responsible to handle the request of get a specific object.
// It returns a gin handler function that handles all the logic behind the http request.
// It uses an s3fetcher interface that is responsible of fetching the stored objects (Openxml, CloudJson, HB, RTA, etc.).
// It calls GetStoredObject that is responsible of obtaiing the objects.
func StorageHandler(s3FetcherUsEast1 storage.S3Fetcher, s3FetcherUsWest1 storage.S3Fetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParams := ExtractGinStorageQueryParams(c)

		status, result, err := api.GetStoredObject(queryParams, s3FetcherUsEast1, s3FetcherUsWest1)

		mimeType := "text"

		if err != nil {
			c.Data(status, mimeType, []byte(err.Error()))
		} else {
			c.Data(status, mimeType, result)
		}
	}
}

// SubscriptionHandler is the responsible to handle requests to retrieve data from Subscriptions
func SubscriptionHandler(fetcher db.PrinterSubscriptionFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryparams := ExtractGinPrinterQueryParams(c)
		status, result, err := api.GetPrinterSubscriptions(queryparams, fetcher)

		if err != nil {
			c.JSON(status, err.Error())
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

// InitRouter initialize a gin router with all the routes for the different endpoints, request types and functions
// that are responsible of handling each request to specific endpoints.
func InitRouter(s3FetcherUsEast1 storage.S3Fetcher, s3FetcherUsWest1 storage.S3Fetcher,
	xmlsFetcher datafetcher.DataFetcher, cloudJsonsFetcher datafetcher.DataFetcher,
	heartbeatsFetcher datafetcher.DataFetcher, rtaFetcher datafetcher.DataFetcher,
	printerSubscriptionFetcher db.PrinterSubscriptionFetcher) *gin.Engine {

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"OPTIONS", "GET", "PUT", "PATCH"},
		AllowHeaders:     []string{"access-control-allow-origin, access-control-allow-headers, Content-Type, x-api-key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWildcard:    true,
		MaxAge:           50 * time.Hour,
	}))

	router.GET(configs.StorageObjectPath, StorageHandler(s3FetcherUsEast1, s3FetcherUsWest1))
	router.GET(configs.OpenXMLPath, Handler(xmlsFetcher))
	router.GET(configs.CloudJsonPath, Handler(cloudJsonsFetcher))
	router.GET(configs.HeartbeatPath, Handler(heartbeatsFetcher))
	router.GET(configs.RTAPath, Handler(rtaFetcher))
	router.GET(configs.SubscriptionsPath, SubscriptionHandler(printerSubscriptionFetcher))

	return router
}
