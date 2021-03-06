package queryparams_test

import (
	"bitbucket.org/aldoft/printer-timeline-backend/app/internal/configs"
	"strconv"
	"time"

	. "bitbucket.org/aldoft/printer-timeline-backend/app/internal/queryparams"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const maxTimeDiffInMinutes = 60

var _ = Describe("Query Parameters controller", func() {
	Describe("Extract Time Range from query parameters", func() {

		Context("Request query parameters don't contains any parameter", func() {
			It("returns missing time range type error", func() {
				queryParams := map[string]string{}
				_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
				Expect(err).To(Equal(ErrorQueryStringMissingTimeRangeType))
			})
		})

		Context("Request Query parameters don't contain time range type but contains other parameters", func() {
			It("returns missing time range type error", func() {
				queryParams := map[string]string{
					configs.ProductNumberQueryParam: "L2E27A",
					configs.SerialNumberQueryParam:  "SG59L1Q005",
				}
				_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
				Expect(err).To(Equal(ErrorQueryStringMissingTimeRangeType))
			})
		})

		Context("Request Query parameters time range type is empty", func() {
			It("returns missing time range type error", func() {
				queryParams := map[string]string{
					configs.ProductNumberQueryParam: "L2E27A",
					configs.SerialNumberQueryParam:  "SG59L1Q005",
					configs.TimeTypeQueryParam:      "",
				}
				_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
				Expect(err).To(Equal(ErrorQueryStringMissingTimeRangeType))
			})
		})

		Context("Request Query parameters contain time range type but is not 'relative' or 'absolute'", func() {
			It("returns query string unsupported time range type error", func() {
				queryParams := map[string]string{
					configs.TimeTypeQueryParam: "always",
				}
				_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
				Expect(err).To(Equal(ErrorQueryStringUnsupportedTimeRangeType))
			})
		})

		Context("Request Query parameters time range type is not supported", func() {
			It("returns unsupported time range type error", func() {
				start := strconv.FormatInt(time.Now().Unix(), 10)
				end := strconv.FormatInt(time.Now().Add(-time.Minute*20).Unix(), 10)
				queryParams := map[string]string{
					configs.TimeTypeQueryParam:  "invented_time_type",
					configs.StartTimeQueryParam: start,
					configs.EndTimeQueryParam:   end,
				}
				_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
				Expect(err).To(Equal(ErrorQueryStringUnsupportedTimeRangeType))
			})
		})

		Context("When time range is relative", func() {
			Context("and start time is present", func() {
				It("returns start time should not appear error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "relative",
						configs.StartTimeQueryParam: "1590084529",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringStartTimeAppears))
				})
			})

			Context("and end time is present", func() {
				It("returns end time should not appear error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam: "relative",
						configs.EndTimeQueryParam:  "1590084529",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringEndTimeAppears))
				})
			})

			Context("and offset units is not present", func() {
				It("returns missing offset units error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam: "relative",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringMissingOffsetUnits))
				})
			})

			Context("and offset units is unsupported (days)", func() {
				It("returns unsupported offset units error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "days",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringUnsupportedOffsetUnits))
				})
			})

			Context("and offset value is missing", func() {
				It("returns missing offset value error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "minutes",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringMissingOffsetValue))
				})
			})

			Context("and offset_units is ok but offset value is not a number", func() {
				It("returns unsupported offset value error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "seconds",
						configs.OffsetValueQueryParam: "Golang",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringUnsupportedOffsetValue))
				})
			})

			Context("and offset units is 'minutes' but offset value is too big", func() {
				It("returns unsupported offset value error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "minutes",
						configs.OffsetValueQueryParam: "61",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringUnsupportedOffsetValue))
				})
			})

			Context("and offset units is minutes but offset value is negative", func() {
				It("returns unsupported offset value error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "minutes",
						configs.OffsetValueQueryParam: "-61",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringUnsupportedOffsetValue))
				})
			})

			Context("and offset units is minutes but offset value is zero", func() {
				It("returns unsupported offset value error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "minutes",
						configs.OffsetValueQueryParam: "0",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringUnsupportedOffsetValue))
				})
			})

			Context("and offset units is minutes and offset value is ok", func() {
				It("returns no error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "minutes",
						configs.OffsetValueQueryParam: "5",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(BeNil())
				})
			})

			Context("and offset units is seconds but offset value is too big", func() {
				It("returns unsupported offset value error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "seconds",
						configs.OffsetValueQueryParam: "36001",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringUnsupportedOffsetValue))
				})
			})

			Context("and offset units is seconds but offset value is negative", func() {
				It("returns unsupported offset value error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "seconds",
						configs.OffsetValueQueryParam: "-1",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringUnsupportedOffsetValue))
				})
			})

			Context("and query params ok", func() {
				It("returns correct startTime and endTime based on query params", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "minutes",
						configs.OffsetValueQueryParam: "5",
					}

					offsetValue, _ := strconv.Atoi(queryParams[configs.OffsetValueQueryParam])

					duration := -1 * time.Minute * time.Duration(offsetValue)

					expectedEndTime := time.Now().UTC()
					expectedStartTime := expectedEndTime.Add(duration)

					startTime, endTime, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)

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

			Context("and query params ok", func() {
				It("returns correct startTime and endTime based on query params", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "minutes",
						configs.OffsetValueQueryParam: "45",
					}

					offsetValue, _ := strconv.Atoi(queryParams[configs.OffsetValueQueryParam])

					duration := -1 * time.Minute * time.Duration(offsetValue)

					expectedEndTime := time.Now().UTC()
					expectedStartTime := expectedEndTime.Add(duration)

					startTime, endTime, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)

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

			Context("and query params ok", func() {
				It("returns correct startTime and endTime based on query params", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "minutes",
						configs.OffsetValueQueryParam: "27",
					}

					offsetValue, _ := strconv.Atoi(queryParams[configs.OffsetValueQueryParam])

					duration := -1 * time.Minute * time.Duration(offsetValue)

					expectedEndTime := time.Now().UTC()
					expectedStartTime := expectedEndTime.Add(duration)

					startTime, endTime, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)

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

			Context("and query params ok", func() {
				It("returns correct startTime and endTime based on query params", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "seconds",
						configs.OffsetValueQueryParam: "30",
					}

					offsetValue, _ := strconv.Atoi(queryParams[configs.OffsetValueQueryParam])
					duration := -1 * time.Second * time.Duration(offsetValue)

					expectedEndTime := time.Now().UTC()
					expectedStartTime := expectedEndTime.Add(duration)

					startTime, endTime, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)

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
					Expect(endTime.Second()).To(Equal(expectedEndTime.Second()))

				})
			})

			Context("and query params ok", func() {
				It("returns correct startTime and endTime based on query params", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:    "relative",
						configs.OffsetUnitsQueryParam: "seconds",
						configs.OffsetValueQueryParam: "50",
					}

					offsetValue, _ := strconv.Atoi(queryParams[configs.OffsetValueQueryParam])
					duration := -1 * time.Second * time.Duration(offsetValue)

					expectedEndTime := time.Now().UTC()
					expectedStartTime := expectedEndTime.Add(duration)

					startTime, endTime, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)

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
					Expect(endTime.Second()).To(Equal(expectedEndTime.Second()))

				})
			})
		})

		Context("When time range is absolute (in UTC)", func() {
			Context("and start time is not present", func() {
				It("returns missing start time error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam: "absolute",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringMissingStartTime))
				})
			})

			Context("and start time is empty", func() {
				It("returns missing start time error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringMissingStartTime))
				})
			})

			Context("and start time has wrong value (is a word)", func() {
				It("returns unsupported start_time error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "This_is_invalid",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringUnsupportedStartTime))
				})
			})

			Context("and start time has wrong value (is a string representing a float)", func() {
				It("returns unsupported start_time error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "6.6",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringUnsupportedStartTime))
				})
			})

			Context("and start time is ok but end time is missing", func() {
				It("returns missing end time error", func() {
					nowEpoch := strconv.FormatInt(time.Now().Unix(), 10)
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: nowEpoch,
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringMissingEndTime))
				})
			})

			Context("and start time is ok but end time is empty", func() {
				It("returns missing end time error", func() {
					now := strconv.FormatInt(time.Now().Unix(), 10)
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: now,
						configs.EndTimeQueryParam:   "",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringMissingEndTime))
				})
			})

			Context("and start time is ok but end time has wrong value (is a word)", func() {
				It("returns unsupported end time error", func() {
					now := strconv.FormatInt(time.Now().Unix(), 10)
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: now,
						configs.EndTimeQueryParam:   "Software",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringUnsupportedEndTime))
				})
			})

			Context("and start time is ok but end time has wrong value (is a float)", func() {
				It("returns unsupported end time error", func() {
					now := strconv.FormatInt(time.Now().Unix(), 10)
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: now,
						configs.EndTimeQueryParam:   "159008452.9",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringUnsupportedEndTime))
				})
			})

			Context("but difference between start time and end time is more than one hour", func() {
				It("returns query string time difference is too big error", func() {
					start := strconv.FormatInt(time.Now().Add(-time.Minute*70).Unix(), 10)
					end := strconv.FormatInt(time.Now().Unix(), 10)
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: start,
						configs.EndTimeQueryParam:   end,
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringTimeDifferenceTooBig))
				})
			})

			Context("but difference between start time and end time is more than one hour", func() {
				It("returns query string time difference is too big error", func() {
					start := strconv.FormatInt(time.Now().Add(-time.Hour*2).Unix(), 10)
					end := strconv.FormatInt(time.Now().Add(time.Minute*3).Unix(), 10)
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: start,
						configs.EndTimeQueryParam:   end,
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringTimeDifferenceTooBig))
				})
			})

			Context("but difference between start time is ok", func() {
				It("returns no error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "1590753903",
						configs.EndTimeQueryParam:   "1590755043",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(BeNil())
				})
			})

			Context("and difference between start time and end time is ok", func() {
				It("returns no error", func() {
					start := strconv.FormatInt(time.Now().Add(-time.Minute*30).Unix(), 10)
					end := strconv.FormatInt(time.Now().Unix(), 10)
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: start,
						configs.EndTimeQueryParam:   end,
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(BeNil())
				})
			})

			Context("and difference between start time and end time is ok", func() {
				It("returns no error", func() {
					start := strconv.FormatInt(time.Now().Add(-time.Minute*60).Unix(), 10)
					end := strconv.FormatInt(time.Now().Unix(), 10)
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: start,
						configs.EndTimeQueryParam:   end,
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(BeNil())
				})
			})

			Context("and difference between start time and end time is ok (30 min)", func() {
				It("returns no error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "1590753642",
						configs.EndTimeQueryParam:   "1590755442",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(BeNil())
				})
			})

			Context("and difference between start time and end time is ok (1 hour)", func() {
				It("returns no error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "1590751918",
						configs.EndTimeQueryParam:   "1590755518",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(BeNil())
				})
			})

			Context("and difference between start time and end time is ok (1 hour)", func() {
				It("returns no error and correct start_time and end_time", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "1590751918",
						configs.EndTimeQueryParam:   "1590755518",
					}
					startTime, endTime, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)

					Expect(err).To(BeNil())

					Expect(startTime.Year()).To(Equal(2020))
					Expect(startTime.Month().String()).To(Equal("May"))
					Expect(startTime.Day()).To(Equal(29))
					Expect(startTime.Hour()).To(Equal(11))
					Expect(startTime.Minute()).To(Equal(31))
					Expect(startTime.Second()).To(Equal(58))

					Expect(endTime.Year()).To(Equal(2020))
					Expect(endTime.Month().String()).To(Equal("May"))
					Expect(endTime.Day()).To(Equal(29))
					Expect(endTime.Hour()).To(Equal(12))
					Expect(endTime.Minute()).To(Equal(31))
					Expect(endTime.Second()).To(Equal(58))
				})
			})

			Context("and difference between start time and end time is ok (30 min)", func() {
				It("returns no error and correct start_time and end_time", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "1589968858",
						configs.EndTimeQueryParam:   "1589970658",
					}
					startTime, endTime, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)

					Expect(err).To(BeNil())

					Expect(startTime.Year()).To(Equal(2020))
					Expect(startTime.Month().String()).To(Equal("May"))
					Expect(startTime.Day()).To(Equal(20))
					Expect(startTime.Hour()).To(Equal(10))
					Expect(startTime.Minute()).To(Equal(0))
					Expect(startTime.Second()).To(Equal(58))

					Expect(endTime.Year()).To(Equal(2020))
					Expect(endTime.Month().String()).To(Equal("May"))
					Expect(endTime.Day()).To(Equal(20))
					Expect(endTime.Hour()).To(Equal(10))
					Expect(endTime.Minute()).To(Equal(30))
					Expect(endTime.Second()).To(Equal(58))
				})
			})

			Context("and difference between start time and end time is ok (45 min)", func() {
				It("returns no error and correct start_time and end_time", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "1590825058",
						configs.EndTimeQueryParam:   "1590827758",
					}
					startTime, endTime, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)

					Expect(err).To(BeNil())

					Expect(startTime.Year()).To(Equal(2020))
					Expect(startTime.Month().String()).To(Equal("May"))
					Expect(startTime.Day()).To(Equal(30))
					Expect(startTime.Hour()).To(Equal(7))
					Expect(startTime.Minute()).To(Equal(50))
					Expect(startTime.Second()).To(Equal(58))

					Expect(endTime.Year()).To(Equal(2020))
					Expect(endTime.Month().String()).To(Equal("May"))
					Expect(endTime.Day()).To(Equal(30))
					Expect(endTime.Hour()).To(Equal(8))
					Expect(endTime.Minute()).To(Equal(35))
					Expect(endTime.Second()).To(Equal(58))
				})
			})

			Context("and difference between start time and end time is ok (20 seconds)", func() {
				It("returns no error and correct start_time and end_time", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "1578215700",
						configs.EndTimeQueryParam:   "1578215720",
					}
					startTime, endTime, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)

					Expect(err).To(BeNil())

					Expect(startTime.Year()).To(Equal(2020))
					Expect(startTime.Month().String()).To(Equal("January"))
					Expect(startTime.Day()).To(Equal(5))
					Expect(startTime.Hour()).To(Equal(9))
					Expect(startTime.Minute()).To(Equal(15))
					Expect(startTime.Second()).To(Equal(00))

					Expect(endTime.Year()).To(Equal(2020))
					Expect(endTime.Month().String()).To(Equal("January"))
					Expect(endTime.Day()).To(Equal(5))
					Expect(endTime.Hour()).To(Equal(9))
					Expect(endTime.Minute()).To(Equal(15))
					Expect(endTime.Second()).To(Equal(20))
				})
			})

			Context("but difference between start time is more than one hour", func() {
				It("returns query string time difference is too big error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "1590751203",
						configs.EndTimeQueryParam:   "1590755043",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringTimeDifferenceTooBig))
				})
			})

			Context("but difference between start time is more than one day", func() {
				It("returns query string time difference is too big error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "1589282403",
						configs.EndTimeQueryParam:   "1590755043",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringTimeDifferenceTooBig))
				})
			})

			Context("and difference between start time and end time is too big", func() {
				It("returns query string time difference is too big error", func() {
					start := strconv.FormatInt(time.Now().Add(-time.Minute*61).Unix(), 10)
					end := strconv.FormatInt(time.Now().Unix(), 10)
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: start,
						configs.EndTimeQueryParam:   end,
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringTimeDifferenceTooBig))
				})
			})

			Context("but end time is previous than start time", func() {
				It("returns query string end time is previous than start time error", func() {
					start := strconv.FormatInt(time.Now().Unix(), 10)
					end := strconv.FormatInt(time.Now().Add(-time.Minute*20).Unix(), 10)
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: start,
						configs.EndTimeQueryParam:   end,
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringEndTimePreviousThanStartTime))
				})
			})

			Context("but end time is previous than start time", func() {
				It("returns query string end time is previous than start time error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "1590766655",
						configs.EndTimeQueryParam:   "1590766643",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringEndTimePreviousThanStartTime))
				})
			})

			Context("but end time is previous than start time", func() {
				It("returns query string end time is previous than start time error", func() {
					queryParams := map[string]string{
						configs.TimeTypeQueryParam:  "absolute",
						configs.StartTimeQueryParam: "1590766681",
						configs.EndTimeQueryParam:   "1590680281",
					}
					_, _, err := ExtractTimeRange(queryParams, maxTimeDiffInMinutes)
					Expect(err).To(Equal(ErrorQueryStringEndTimePreviousThanStartTime))
				})
			})
		})

	})

	Describe("Extract Printer number and Serial number from query parameters", func() {
		Context("When neither Product number nor Serial number is present in query parameters", func() {
			It("returns no error", func() {
				queryParams := map[string]string{}
				_, _, err := ExtractPrinterInfo(queryParams)

				Expect(err).To(BeNil())

			})
		})

		Context("When Product number is present but Serial Number is missing", func() {
			It("returns no error", func() {
				queryParams := map[string]string{
					configs.ProductNumberQueryParam: "CZ056A",
				}
				_, _, err := ExtractPrinterInfo(queryParams)

				Expect(err).To(BeNil())

			})
		})

		Context("When Serial number is present but Product number is missing", func() {
			It("returns no error and correct Pn and Sn", func() {
				queryParams := map[string]string{
					configs.SerialNumberQueryParam: "SG4491P001",
				}
				_, _, err := ExtractPrinterInfo(queryParams)

				Expect(err).To(Equal(ErrorQueryStringPnSn))

			})
		})

		Context("When Serial number is present but Product number is missing", func() {
			It("returns no error and correct Pn and Sn", func() {
				queryParams := map[string]string{
					configs.ProductNumberQueryParam: "CZ056A",
					configs.SerialNumberQueryParam:  "SG4491P001",
				}
				productNumber, serialNumber, err := ExtractPrinterInfo(queryParams)

				Expect(productNumber).To(Equal("CZ056A"))
				Expect(serialNumber).To(Equal("SG4491P001"))
				Expect(err).To(BeNil())

			})
		})
	})

	Describe("Extract Printer number and Serial number from query parameters", func() {
		Context("When neither Product number nor Serial number is present in query parameters", func() {
			It("returns no error", func() {
				queryParams := map[string]string{}
				_, _, err := ExtractPrinterInfo(queryParams)

				Expect(err).To(BeNil())

			})
		})

		Context("When Product number is present but Serial Number is missing", func() {
			It("returns no error", func() {
				queryParams := map[string]string{
					configs.ProductNumberQueryParam: "CZ056A",
				}
				_, _, err := ExtractPrinterInfo(queryParams)

				Expect(err).To(BeNil())

			})
		})

		Context("When Serial number is present but Product number is missing", func() {
			It("returns no error and correct Pn and Sn", func() {
				queryParams := map[string]string{
					configs.SerialNumberQueryParam: "SG4491P001",
				}
				_, _, err := ExtractPrinterInfo(queryParams)

				Expect(err).To(Equal(ErrorQueryStringPnSn))

			})
		})

		Context("When Serial number is present but Product number is missing", func() {
			It("returns no error and correct Pn and Sn", func() {
				queryParams := map[string]string{
					configs.ProductNumberQueryParam: "CZ056A",
					configs.SerialNumberQueryParam:  "SG4491P001",
				}
				productNumber, serialNumber, err := ExtractPrinterInfo(queryParams)

				Expect(productNumber).To(Equal("CZ056A"))
				Expect(serialNumber).To(Equal("SG4491P001"))
				Expect(err).To(BeNil())
			})
		})
	})

})
