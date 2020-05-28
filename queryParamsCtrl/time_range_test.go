package queryParamsCtrl_test

import (
	"bitbucket.org/aldoft/printer-timeline-backend/errors"
	"bitbucket.org/aldoft/printer-timeline-backend/queryParamsCtrl"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
	"time"
)

var _ = Describe("Time range controller", func() {
	Describe("Extract Time Range from query parameters", func() {

		Context("Request query parameters not contain any parameter", func() {
			It("returns missing time range type error", func() {
				queryParams := map[string]string{
				}
				 _, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringMissingTimeRangeType))
			})
		})

		Context("Request Query parameters not contain time range type but contain other parameters", func() {
			It("returns missing time range type error", func() {
				queryParams := map[string]string{
					"pn" : "L2E27A",
					"sn": "SG59L1Q005",
				}
				 _, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringMissingTimeRangeType))
			})
		})

		Context("Request Query parameters time range type is empty", func() {
			It("returns missing time range type error", func() {
				queryParams := map[string]string{
					"pn" : "L2E27A",
					"sn": "SG59L1Q005",
					"time_type": "",
				}
				 _, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringMissingTimeRangeType))
			})
		})

		Context("Request Query parameters contain time range type but are not 'relative' or 'absolute'", func() {
			It("returns query string unsupported time range type error", func() {
				queryParams := map[string]string{
					"time_type": "always",
				}
				 _, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringUnsupportedTimeRangeType))
			})
		})

		Context("Request Query parameters time range type is relative and start time is present", func() {
			It("returns start time should not appear error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"start_time": "1590084529",
				}
				 _, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringStartTimeAppears))
			})
		})

		Context("Request Query parameters time range type is relative and end time is present", func() {
			It("returns end time should not appear error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"end_time": "1590084529",
				}
				 _, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringEndTimeAppears))
			})
		})

		Context("Request Query parameters time range type is relative and offset units is not present", func() {
			It("returns missing offset units error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
				}
				 _, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringMissingOffsetUnits))
			})
		})

		Context("Request Query parameters time range type is relative and offset units is unsupported", func() {
			It("returns unsupported offset units error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "days", //days are not supported for now
				}
				 _, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringUnsupportedOffsetUnits))
			})
		})

		Context("Request Query parameters time range type is relative and offset value is missing", func() {
			It("returns missing offset value error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "minutes",
				}
				 _, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringMissingOffsetValue))
			})
		})

		Context("Request Query parameters time range type and offset_units are ok but offset value is not a number", func() {
			It("returns unsupported offset value error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "seconds",
					"offset_value": "Golang",
				}
				 _, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringUnsupportedOffsetValue))
			})
		})

		Context("Request Query parameters time range type is relative, offset units is minutes but offset value is too big", func() {
			It("returns unsupported offset value error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "minutes",
					"offset_value": "61",
				}
				 _, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringUnsupportedOffsetValue))
			})
		})

		Context("Request Query parameters time range type is relative, offset units is minutes but offset value is negative", func() {
			It("returns unsupported offset value error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "minutes",
					"offset_value": "-61",
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringUnsupportedOffsetValue))
			})
		})

		Context("Request Query parameters time range type is relative, offset units is minutes but offset value is zero", func() {
			It("returns unsupported offset value error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "minutes",
					"offset_value": "0",
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringUnsupportedOffsetValue))
			})
		})

		Context("Request Query parameters time range type is relative, offset units is minutes and offset value is ok", func() {
			It("returns no error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "minutes",
					"offset_value": "5",
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
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
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringUnsupportedOffsetValue))
			})
		})

		Context("Request Query parameters time range type is relative, offset units is seconds but offset value is negative", func() {
			It("returns unsupported offset value error", func() {
				queryParams := map[string]string{
					"time_type": "relative",
					"offset_units": "seconds",
					"offset_value": "-1",
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringUnsupportedOffsetValue))
			})
		})

		Context("Request Query parameters time range type is absolute and start time is not present", func() {
			It("returns missing start time error", func() {
				queryParams := map[string]string{
					"time_type": "absolute",
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringMissingStartTime))
			})
		})

		Context("Request Query parameters time range type is absolute and start time is empty", func() {
			It("returns missing start time error", func() {
				queryParams := map[string]string{
					"time_type": "absolute",
					"start_time": "",
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringMissingStartTime))
			})
		})

		Context("Request Query parameters time range type is absolute and start time has wrong value (is a word)", func() {
			It("returns unsupported start_time error", func() {
				queryParams := map[string]string{
					"time_type": "absolute",
					"start_time": "This_is_invalid",
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringUnsupportedStartTime))
			})
		})

		Context("Request Query parameters time range type is absolute and start time has wrong value (is a float)", func() {
			It("returns unsupported start_time error", func() {
				queryParams := map[string]string{
					"time_type": "absolute",
					"start_time": "6.6",
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringUnsupportedStartTime))
			})
		})

		Context("Request Query parameters time range type is absolute, start time is ok but end time is missing", func() {
			It("returns missing end time error", func() {
				nowEpoch := strconv.FormatInt(time.Now().Unix(), 10)
				queryParams := map[string]string{
					"time_type":  "absolute",
					"start_time": nowEpoch,
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringMissingEndTime))
			})
		})

		Context("Request Query parameters time range type is absolute, start time is ok but end time is empty", func() {
			It("returns missing end time error", func() {
				now := strconv.FormatInt(time.Now().Unix(), 10)
				queryParams := map[string]string{
					"time_type":  "absolute",
					"start_time": now,
					"end_time":   "",
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringMissingEndTime))
			})
		})

		Context("Request Query parameters time range type is absolute, start time is ok but end time has wrong value (is a word)", func() {
			It("returns unsupported end time error", func() {
				now := strconv.FormatInt(time.Now().Unix(), 10)
				queryParams := map[string]string{
					"time_type":  "absolute",
					"start_time": now,
					"end_time":   "Software",
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringUnsupportedEndTime))
			})
		})

		Context("Request Query parameters time range type is absolute, start time is ok but end time has wrong value (is a float)", func() {
			It("returns unsupported end time error", func() {
				now := strconv.FormatInt(time.Now().Unix(), 10)
				queryParams := map[string]string{
					"time_type":  "absolute",
					"start_time": now,
					"end_time":   "159008452.9",
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringUnsupportedEndTime))
			})
		})

		Context("Request Query parameters time range type is absolute but difference between start time and end time is more than one hour", func() {
			It("returns query string time difference is too big error", func() {
				start := strconv.FormatInt(time.Now().Add(-time.Minute * 70).Unix(), 10)
				end := strconv.FormatInt(time.Now().Unix(), 10)
				queryParams := map[string]string{
					"time_type":  "absolute",
					"start_time": start,
					"end_time":   end,
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringTimeDifferenceTooBig))
			})
		})

		Context("Request Query parameters time range type is absolute but difference between start time and end time is more than one hour", func() {
			It("returns query string time difference is too big error", func() {
				start := strconv.FormatInt(time.Now().Add(-time.Hour * 2).Unix(), 10)
				end := strconv.FormatInt(time.Now().Add(time.Minute * 3).Unix(), 10)
				queryParams := map[string]string{
					"time_type":  "absolute",
					"start_time": start,
					"end_time":   end,
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringTimeDifferenceTooBig))
			})
		})

		Context("Request Query parameters time range type is absolute and difference between start time and end time is ok", func() {
			It("returns no error", func() {
				start := strconv.FormatInt(time.Now().Add(-time.Minute * 30).Unix(), 10)
				end := strconv.FormatInt(time.Now().Unix(), 10)
				queryParams := map[string]string{
					"time_type":  "absolute",
					"start_time": start,
					"end_time":   end,
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(BeNil())
			})
		})

		Context("Request Query parameters time range type is absolute and difference between start time and end time is ok", func() {
			It("returns no error", func() {
				start := strconv.FormatInt(time.Now().Add(-time.Minute * 60).Unix(), 10)
				end := strconv.FormatInt(time.Now().Unix(), 10)
				queryParams := map[string]string{
					"time_type":  "absolute",
					"start_time": start,
					"end_time":   end,
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(BeNil())
			})
		})

		Context("Request Query parameters time range type is absolute and difference between start time and end time is too big", func() {
			It("returns query string time difference is too big error", func() {
				start := strconv.FormatInt(time.Now().Add(-time.Minute * 61).Unix(), 10)
				end := strconv.FormatInt(time.Now().Unix(), 10)
				queryParams := map[string]string{
					"time_type":  "absolute",
					"start_time": start,
					"end_time":   end,
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringTimeDifferenceTooBig))
			})
		})

		Context("Request Query parameters time range type is absolute but end time is previous than start time", func() {
			It("returns query string end time is previous than start time error", func() {
				start := strconv.FormatInt(time.Now().Unix(), 10)
				end := strconv.FormatInt(time.Now().Add(-time.Minute * 20).Unix(), 10)
				queryParams := map[string]string{
					"time_type":  "absolute",
					"start_time": start,
					"end_time":   end,
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringEndTimePreviousThanStartTime))
			})
		})

		Context("Request Query parameters time range type is not supported", func() {
			It("returns unsupported time range type error", func() {
				start := strconv.FormatInt(time.Now().Unix(), 10)
				end := strconv.FormatInt(time.Now().Add(-time.Minute * 20).Unix(), 10)
				queryParams := map[string]string{
					"time_type":  "invented_time_type",
					"start_time": start,
					"end_time":   end,
				}
				_, _, err := queryParamsCtrl.ExtractTimeRange(queryParams)
				Expect(err).To(Equal(errors.QueryStringUnsupportedTimeRangeType))
			})
		})

		Context("Relative time and query params ok", func() {
			It("returns correct startTime and endTime based on query params", func() {
				queryParams := map[string]string{
					"time_type":  "relative",
					"offset_units": "minutes",
					"offset_value":   "5",
				}
				expectedEndTime := time.Now()
				startTime, endTime, err := queryParamsCtrl.ExtractTimeRange(queryParams)

				offsetValue, _ := strconv.Atoi(queryParams["offset_value"])

				duration := -1 * time.Minute * time.Duration(offsetValue)

				expectedStartTime := time.Now().Add(duration)

				Expect(err).To(BeNil())

				Expect(startTime.Year()).To(Equal(expectedStartTime.Year()))
				Expect(startTime.Month()).To(Equal(expectedStartTime.Month()))
				Expect(startTime.Day()).To(Equal(expectedStartTime.Day()))
				Expect(startTime.Hour()).To(Equal(expectedStartTime.Hour()))
				Expect(startTime.Minute()).To(Equal(expectedStartTime.Minute()))
				Expect(startTime.Second()).To(Equal(expectedStartTime.Second()))

				Expect(endTime.Year()).To(Equal(expectedEndTime.Year()))
				Expect(endTime.Month()).To(Equal(expectedEndTime.Month()))
				Expect(endTime.Day()).To(Equal(expectedEndTime.Day()))
				Expect(endTime.Hour()).To(Equal(expectedEndTime.Hour()))
				Expect(endTime.Minute()).To(Equal(expectedEndTime.Minute()))
				Expect(endTime.Second()).To(Equal(expectedStartTime.Second()))

			})
		})
	})
})
