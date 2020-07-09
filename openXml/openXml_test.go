package openXml_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bitbucket.org/aldoft/printer-timeline-backend/openXml"
)

var _ = Describe("OpenXml service", func() {
	Describe("GetUploadedOpenXmls", func() {
		Context("When one of the HTTP query parameters is not correct", func() {
			It("returns the appropiate error", func() {
				status, response := openXml.GetUploadedOpenXmls(svc, queryParameters)
				
				Expect(err).To(Equal(errors.QueryStringMissingTimeRangeType))
			})
		})
	}
})
