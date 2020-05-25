package openXml_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bitbucket.org/aldoft/printer-timeline-backend/api/common"
	"bitbucket.org/aldoft/printer-timeline-backend/api/openXml"
)

var _ = Describe("OpenXml", func() {
	Describe("Prepare insights query parameters", func() {

		Context("Request query parameters not contain any parameter", func() {
			It("returns missing time range type error", func() {
				queryParams := map[string]string{
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringMissingTimeRangeTypeError))
			})
		})

		Context("Request Query parameters not contain time range type but contain other parameters", func() {
			It("returns missing time range type error", func() {
				queryParams := map[string]string{
					"pn" : "L2E27A",
					"sn": "SG59L1Q005",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringMissingTimeRangeTypeError))
			})
		})

		Context("Request Query parameters time range type is empty", func() {
			It("returns missing time range type error", func() {
				queryParams := map[string]string{
					"pn" : "L2E27A",
					"sn": "SG59L1Q005",
					"time_type": "",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringMissingTimeRangeTypeError))
			})
		})

		Context("Request Query parameters contain time range type but are not 'relative' or 'absolute'", func() {
			It("returns query string unsupported time range type error", func() {
				queryParams := map[string]string{
					"time_type": "always",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringUnsupportedTimeRangeTypeError))
			})
		})

		Context("Request Query parameters time range type is relative and start time is present", func() {
			It("returns start time should not appear error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"start_time": "1590084529",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringPresentStartTimeError))
			})
		})

		Context("Request Query parameters time range type is relative and end time is present", func() {
			It("returns end time should not appear error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"end_time": "1590084529",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringPresentEndTimeError))
			})
		})

		Context("Request Query parameters time range type is relative and offset units is not present", func() {
			It("returns missing offset units error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringMissingOffsetUnitsError))
			})
		})

		Context("Request Query parameters time range type is relative and offset units is unsupported", func() {
			It("returns unsupported offset units error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "days",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringUnsupportedOffsetUnitsError))
			})
		})

		Context("Request Query parameters time range type is relative, offset units is ok but offset value is too big", func() {
			It("returns unsupported offset value error because is too big", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "minutes",
					"offset_value": "61",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringUnsupportedOffsetValueError))
			})
		})
	})
})
