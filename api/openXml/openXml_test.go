package openXml_test

import (
	"bitbucket.org/aldoft/printer-timeline-backend/api/common"
	"bitbucket.org/aldoft/printer-timeline-backend/api/openXml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
	"time"
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
				Expect(err).To(Equal(common.QueryStringStartTimeAppearsError))
			})
		})

		Context("Request Query parameters time range type is relative and end time is present", func() {
			It("returns end time should not appear error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"end_time": "1590084529",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringEndTimeAppearsError))
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
					"offset_units": "days", //days are not supported for now
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringUnsupportedOffsetUnitsError))
			})
		})

		Context("Request Query parameters time range type is relative and offset value is missing", func() {
			It("returns missing offset value error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "minutes",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringMissingOffsetValueError))
			})
		})

		Context("Request Query parameters time range type and offset_units are ok but offset value is not a number", func() {
			It("returns unsupported offset value error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "seconds",
					"offset_value": "Golang",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringUnsupportedOffsetValueError))
			})
		})

		Context("Request Query parameters time range type is relative, offset units is minutes but offset value is too big", func() {
			It("returns unsupported offset value error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "minutes",
					"offset_value": "61",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringUnsupportedOffsetValueError))
			})
		})

		Context("Request Query parameters time range type is relative, offset units is minutes but offset value is negative", func() {
			It("returns unsupported offset value error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "minutes",
					"offset_value": "-61",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringUnsupportedOffsetValueError))
			})
		})

		Context("Request Query parameters time range type is relative, offset units is minutes but offset value is zero", func() {
			It("returns unsupported offset value error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "minutes",
					"offset_value": "0",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringUnsupportedOffsetValueError))
			})
		})

		Context("Request Query parameters time range type is relative, offset units is minutes and offset value is ok", func() {
			It("returns no error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "minutes",
					"offset_value": "5",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(BeNil())
			})
		})

		Context("Request Query parameters time range type is relative, offset units is seconds but offset value is too big", func() {
			It("returns unsupported offset value error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "seconds",
					"offset_value": "36001",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringUnsupportedOffsetValueError))
			})
		})

		Context("Request Query parameters time range type is relative, offset units is seconds but offset value is negative", func() {
			It("returns unsupported offset value error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "seconds",
					"offset_value": "-1",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringUnsupportedOffsetValueError))
			})
		})

		Context("Request Query parameters time range type is absolute and start time is not present", func() {
			It("returns missing start time error", func() {
				queryParams := map[string]string{
					"time_type": "absolute",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringMissingStartTimeError))
			})
		})

		Context("Request Query parameters time range type is absolute and start time is empty", func() {
			It("returns missing start time error", func() {
				queryParams := map[string]string{
					"time_type": "absolute",
					"start_time": "",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringMissingStartTimeError))
			})
		})

		Context("Request Query parameters time range type is absolute and start time has wrong value (is a word)", func() {
			It("returns unsupported start_time", func() {
				queryParams := map[string]string{
					"time_type": "absolute",
					"start_time": "This_is_invalid",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringUnsupportedStartTimeError))
			})
		})

		Context("Request Query parameters time range type is absolute, start time is ok but end_time is missing", func() {
			It("returns missing end time error", func() {
				now := strconv.FormatInt(time.Now().Unix(), 10)
				queryParams := map[string]string{
					"time_type": "absolute",
					"start_time": now,
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringMissingEndTimeError))
			})
		})

		Context("Request Query parameters time range type is absolute, start time is ok but end_time is empty", func() {
			It("returns missing end time error", func() {
				now := strconv.FormatInt(time.Now().Unix(), 10)
				queryParams := map[string]string{
					"time_type": "absolute",
					"start_time": now,
					"end_time": "",
				}
				_, _, _, err := openXml.PrepareInsightsQueryParameters(queryParams)
				Expect(err).To(Equal(common.QueryStringMissingEndTimeError))
			})
		})
	})
})
