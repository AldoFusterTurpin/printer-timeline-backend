package api_test

import (
	"errors"

	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/api"
	myErrors "bitbucket.org/aldoft/printer-timeline-backend/app/internal/errors"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/openXml/mocks"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Api", func() {
	Describe("OpenXMLHandler", func() {
		Context("When there is an error in query parameters", func() {
			It("returns the appropiate result, status and error (correctly propagating the error)", func() {
				mockCtrl := gomock.NewController(GinkgoT())
				defer mockCtrl.Finish()

				mockOpenXmlsFetcher := mocks.NewMockOpenXmlsFetcher(mockCtrl)

				mockOpenXmlsFetcher.EXPECT().GetUploadedOpenXmls(gomock.Any()).Return(nil, myErrors.QueryStringMissingTimeRangeType).Times(1)

				queryparams := map[string]string{}
				status, result, err := api.GetOpenXmls(queryparams, mockOpenXmlsFetcher)

				Expect(status).To(Equal(400))
				Expect(result).To(BeNil())
				Expect(err).To(Equal(myErrors.QueryStringMissingTimeRangeType))
			})
		})
	})

	Describe("SelectHTTPStatus", func() {
		Context("When the input is QueryStringMissingTimeRangeType error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringMissingTimeRangeType)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringUnsupportedTimeRangeType error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringUnsupportedTimeRangeType)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringStartTimeAppears error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringStartTimeAppears)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringMissingEndTime error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringMissingEndTime)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringEndTimeAppears error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringEndTimeAppears)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringUnsupportedEndTime error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringUnsupportedEndTime)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringMissingOffsetUnits error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringMissingOffsetUnits)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringUnsupportedOffsetUnits error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringUnsupportedOffsetUnits)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringMissingOffsetValue error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringMissingOffsetValue)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringUnsupportedOffsetValue error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringUnsupportedOffsetValue)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringMissingStartTime error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringMissingStartTime)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringUnsupportedStartTime error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringUnsupportedStartTime)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringStartTimeAppears error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringMissingTimeRangeType)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringTimeDifferenceTooBig error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringTimeDifferenceTooBig)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringEndTimePreviousThanStartTime error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringEndTimePreviousThanStartTime)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is QueryStringPnSn error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(myErrors.QueryStringPnSn)
				Expect(status).To(Equal(400))
			})
		})

		Context("When the input is nil error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(nil)
				Expect(status).To(Equal(200))
			})
		})

		Context("When the input is an uknown error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(errors.New("invented error"))
				Expect(status).To(Equal(500))
			})
		})
	})

})
