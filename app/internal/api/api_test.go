package api_test

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/api"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/errors"
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/openXml/mocks"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Api", func() {
	Describe("OpenXMLHandler", func() {
		Context("When in HTTP GET query string time range type is missing", func() {
			It("returns the appropiate error", func() {
				mockCtrl := gomock.NewController(GinkgoT())
				defer mockCtrl.Finish()

				mockOpenXmlsFetcher := mocks.NewMockOpenXmlsFetcher(mockCtrl)

				mockOpenXmlsFetcher.EXPECT().GetUploadedOpenXmls(gomock.Any()).Return(nil, errors.QueryStringMissingTimeRangeType).Times(1)

				queryparams := map[string]string{}
				status, result, err := api.GetOpenXmls(queryparams, mockOpenXmlsFetcher)

				Expect(status).To(Equal(400))
				Expect(result).To(BeNil())
				Expect(err).To(Equal(errors.QueryStringMissingTimeRangeType))
			})
		})
	})
})
