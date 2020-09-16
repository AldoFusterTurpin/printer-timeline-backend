package awslambda_test

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/awslambda"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/datafetcher/mocks"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/queryparams"
	s3Mocks "bitbucket.org/aldoft/printer-timeline-backend/app/internal/storage/mocks"
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
)

//go:generate mockgen -source=../s3storage/s3.go -destination=../s3storage/mocks/s3.go -package=mocks

const (
	invalidRegion = "EU_CENTRAL_1"
	bucketName    = "bucketName"
	objectName    = "objectKey"

	testResponse = `{
  "Results": [
    [
      {
        "Field": "@timestamp",
        "Value": "2020-09-16 09:25:00.347"
      },
      {
        "Field": "fields.ProductNumber",
        "Value": "Y0U23A"
      },
      {
        "Field": "fields.SerialNumber",
        "Value": "MY97F1T00H"
      },
      {
        "Field": "fields.bucket_name",
        "Value": "drp-cloudconnector-to-blacksea"
      },
      {
        "Field": "fields.bucket_region",
        "Value": "US_WEST_1"
      },
      {
        "Field": "fields.key",
        "Value": "Y0U23A!MY97F1T00H-7396e65a-f5f2-4db9-8d02-bd1a7d4c8132"
      },
      {
        "Field": "fields.topic",
        "Value": "json"
      },
      {
        "Field": "fields.metadata.date",
        "Value": "2020-09-16T09:24:59Z"
      },
      {
        "Field": "fields.metadata.xml-generator-object-path",
        "Value": "uploads/raw/Y0U23A!MY97F1T00H_2020_09_16_09_24_57_813"
      },
      {
        "Field": "@ptr",
        "Value": "Cl8KJgoiMTAzNjkxMDEzODI3Oi9hd3MvbGFtYmRhL0FXU1BhcnNlchAFEjUaGAIF8qpFWAAAAACOybG5AAX2HZngAAAHgiABKLP57LHJLjCz9O2xyS44CED3E0jRSFDANxABGAE=n     "
      }
    ]
  ],
  "Statistics": {
    "BytesScanned": 16792,
    "RecordsMatched": 13,
    "RecordsScanned": 52
  },
  "Status": "Complete"
}`
)

var _ = Describe("Handlers test", func() {

	Context("CreateLambdaHandler", func() {

		var mockCtrl *gomock.Controller
		var mockCloudJSONFetcher *mocks.MockDataFetcher
		var mockXMLFetcher *mocks.MockDataFetcher
		var mockHeartbeatFetcher *mocks.MockDataFetcher
		var mockRTAFetcher *mocks.MockDataFetcher

		var mockS3UsEastFetcher, mockS3UsWestFetcher *s3Mocks.MockS3Fetcher
		var eventRequest *events.APIGatewayProxyRequest

		var objectResult *s3.GetObjectOutput
		var results *cloudwatchlogs.GetQueryResultsOutput

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())

			mockS3UsEastFetcher = s3Mocks.NewMockS3Fetcher(mockCtrl)
			mockS3UsWestFetcher = s3Mocks.NewMockS3Fetcher(mockCtrl)
			mockXMLFetcher = mocks.NewMockDataFetcher(mockCtrl)
			mockCloudJSONFetcher = mocks.NewMockDataFetcher(mockCtrl)
			mockHeartbeatFetcher = mocks.NewMockDataFetcher(mockCtrl)
			mockRTAFetcher = mocks.NewMockDataFetcher(mockCtrl)

			objectResult = &s3.GetObjectOutput{
				Body:     ioutil.NopCloser(bytes.NewReader([]byte(`content`))),
				Metadata: nil,
			}
			eventRequest = &events.APIGatewayProxyRequest{
				Path:                  "",
				RequestContext:        events.APIGatewayProxyRequestContext{},
				Body:                  "body",
				QueryStringParameters: make(map[string]string),
			}

			err := json.Unmarshal([]byte(testResponse), &results)
			Expect(err).To(BeNil())

		})

		AfterEach(func() {
			defer mockCtrl.Finish()
		})

		It("should call cloudjson fetcher", func() {
			eventRequest.Path = awslambda.CloudJsonPath

			handler := awslambda.CreateLambdaHandler(mockS3UsEastFetcher, mockS3UsWestFetcher, mockXMLFetcher, mockCloudJSONFetcher, mockHeartbeatFetcher, mockRTAFetcher)
			mockCloudJSONFetcher.EXPECT().FetchData(gomock.Any()).Return(results, nil).MinTimes(1)

			_, _ = handler(context.Background(), eventRequest)

		})

		It("should call openXml fetcher", func() {
			eventRequest.Path = awslambda.OpenXMLPath

			handler := awslambda.CreateLambdaHandler(mockS3UsEastFetcher, mockS3UsWestFetcher, mockXMLFetcher, mockCloudJSONFetcher, mockHeartbeatFetcher, mockRTAFetcher)
			mockXMLFetcher.EXPECT().FetchData(gomock.Any()).Return(results, nil).MinTimes(1)

			_, _ = handler(context.Background(), eventRequest)

		})

		It("should call heartbeat fetcher", func() {
			eventRequest.Path = awslambda.HeartbeatPath

			handler := awslambda.CreateLambdaHandler(mockS3UsEastFetcher, mockS3UsWestFetcher, mockXMLFetcher, mockCloudJSONFetcher, mockHeartbeatFetcher, mockRTAFetcher)
			mockHeartbeatFetcher.EXPECT().FetchData(gomock.Any()).Return(results, nil).MinTimes(1)

			_, _ = handler(context.Background(), eventRequest)

		})

		It("should call rta fetcher", func() {
			eventRequest.Path = awslambda.RTAPath

			handler := awslambda.CreateLambdaHandler(mockS3UsEastFetcher, mockS3UsWestFetcher, mockXMLFetcher, mockCloudJSONFetcher, mockHeartbeatFetcher, mockRTAFetcher)
			mockRTAFetcher.EXPECT().FetchData(gomock.Any()).Return(results, nil).MinTimes(1)

			_, _ = handler(context.Background(), eventRequest)

		})

		Context("object tests", func() {
			BeforeEach(func() {
				eventRequest.Path = awslambda.StorageObjectPath

				eventRequest.QueryStringParameters["bucket_name"] = bucketName
				eventRequest.QueryStringParameters["object_key"] = objectName
			})

			It("should call object fetcher of east1 region", func() {
				eventRequest.QueryStringParameters["bucket_region"] = *aws.String(queryparams.UsEast1S3Region)

				handler := awslambda.CreateLambdaHandler(mockS3UsEastFetcher, mockS3UsWestFetcher, mockXMLFetcher, mockCloudJSONFetcher, mockHeartbeatFetcher, mockRTAFetcher)
				mockS3UsEastFetcher.EXPECT().GetObject(gomock.Any()).Return(objectResult, nil).MinTimes(1)

				_, _ = handler(context.Background(), eventRequest)
			})

			It("should call object fetcher of west1 region", func() {
				eventRequest.QueryStringParameters["bucket_region"] = *aws.String(queryparams.UsWest1S3Region)

				handler := awslambda.CreateLambdaHandler(mockS3UsEastFetcher, mockS3UsWestFetcher, mockXMLFetcher, mockCloudJSONFetcher, mockHeartbeatFetcher, mockRTAFetcher)
				mockS3UsWestFetcher.EXPECT().GetObject(gomock.Any()).Return(objectResult, nil).MinTimes(1)

				_, _ = handler(context.Background(), eventRequest)
			})

			It("should not call object fetchers for an invalid region", func() {
				eventRequest.QueryStringParameters["bucket_region"] = *aws.String(invalidRegion)

				handler := awslambda.CreateLambdaHandler(mockS3UsEastFetcher, mockS3UsWestFetcher, mockXMLFetcher, mockCloudJSONFetcher, mockHeartbeatFetcher, mockRTAFetcher)

				resp, _ := handler(context.Background(), eventRequest)

				Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusInternalServerError))
			})

		})

	})

})
