package api_test

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/datafetcher/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"net/http"

	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/api"
	. "bitbucket.org/aldoft/printer-timeline-backend/app/internal/queryparams"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Api", func() {
	Describe("DataFetcher generic", func() {
		Context("When there is an error in query parameters", func() {
			It("returns the appropiate result, status and error (correctly propagating the error)", func() {
				mockCtrl := gomock.NewController(GinkgoT())
				defer mockCtrl.Finish()

				mockDataFetcher := mocks.NewMockDataFetcher(mockCtrl)
				mockDataFetcher.EXPECT().FetchData(gomock.All()).Return(nil, ErrorQueryStringMissingTimeRangeType).Times(1)

				queryparams := map[string]string{}
				status, result, err := api.GetData(queryparams, mockDataFetcher)

				Expect(status).To(Equal(http.StatusBadRequest))
				Expect(result).To(BeNil())
				Expect(err).To(Equal(ErrorQueryStringMissingTimeRangeType))
			})
		})
	})

	Describe("SelectHTTPStatus", func() {
		Context("When the input is QueryStringMissingTimeRangeType error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringMissingTimeRangeType)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringUnsupportedTimeRangeType error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringUnsupportedTimeRangeType)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringStartTimeAppears error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringStartTimeAppears)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringMissingEndTime error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringMissingEndTime)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringEndTimeAppears error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringEndTimeAppears)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringUnsupportedEndTime error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringUnsupportedEndTime)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringMissingOffsetUnits error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringMissingOffsetUnits)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringUnsupportedOffsetUnits error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringUnsupportedOffsetUnits)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringMissingOffsetValue error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringMissingOffsetValue)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringUnsupportedOffsetValue error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringUnsupportedOffsetValue)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringMissingStartTime error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringMissingStartTime)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringUnsupportedStartTime error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringUnsupportedStartTime)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringStartTimeAppears error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringMissingTimeRangeType)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringTimeDifferenceTooBig error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringTimeDifferenceTooBig)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringEndTimePreviousThanStartTime error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringEndTimePreviousThanStartTime)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is QueryStringPnSn error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(ErrorQueryStringPnSn)
				Expect(status).To(Equal(http.StatusBadRequest))
			})
		})

		Context("When the input is nil error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(nil)
				Expect(status).To(Equal(http.StatusOK))
			})
		})

		Context("When the input is an uknown error", func() {
			It("returns the appropiate status", func() {
				status := api.SelectHTTPStatus(errors.New("invented error"))
				Expect(status).To(Equal(http.StatusInternalServerError))
			})
		})
	})

})
