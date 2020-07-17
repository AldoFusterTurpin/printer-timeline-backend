package main_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bitbucket.org/aldoft/printer-timeline-backend/errors"
	"bitbucket.org/aldoft/printer-timeline-backend/mocks"
)

var _ = Describe("Main", func() {
	Describe("OpenXmlHandler function for handle the GET request of uploaded OpneXmls", func() {
		Context("When time range type is missing", func() {
			It("returns the appropiate error", func() {
				mockCtrl := gomock.NewController(GinkgoT())
				defer mockCtrl.Finish()

				//TODO: can't acces functions from main. Should move main/OpenXmlHandler to another package (maybe called apiHandler) but 
				//don't want to mix http status with error types.
				mockOpenXmlsFetcher := mocks.NewMockOpenXmlsFetcher
				resultFunction := OpenXmlHandler(xmlsFetcher)

				mockData.EXPECT().GetUploadedOpenXmls(gomock.Any()).Return(errors.QueryStringMissingTimeRangeType).Times(1)

				Expect(err).To(Equal(errors.QueryStringMissingTimeRangeType))
			})
		})
	})
})
