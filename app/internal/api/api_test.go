package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bitbucket.org/aldoft/printer-timeline-backend/app/internal/api"
)

var _ = Describe("Api", func() {
	Describe("OpenXmlHandler", func() {
		Context("When in HTTP GET query string time range type is missing", func() {
			It("returns the appropiate error", func() {
				mockCtrl := gomock.NewController(GinkgoT())
				defer mockCtrl.Finish()

				mockOpenXmlsFetcher := mocks.NewMockOpenXmlsFetcher
				mockOpenXmlsFetcher.EXPECT().GetUploadedOpenXmls(gomock.Any()).Return(errors.QueryStringMissingTimeRangeType).Times(1)

				resultFunction := OpenXmlHandler(mockOpenXmlsFetcher)

				Expect(err).To(Equal(errors.QueryStringMissingTimeRangeType))
			})
		})
	}
})
