package errors

import "errors"

var (
	QueryStringMissingTimeRangeType     = errors.New("query string missing time range type error")
	QueryStringUnsupportedTimeRangeType = errors.New("query string unsupported time range type error")

	QueryStringStartTimeAppears   = errors.New("query string start time should not appear error")
	QueryStringMissingEndTime     = errors.New("query string missing end time when time range is absolute error")
	QueryStringEndTimeAppears     = errors.New("query string end time should not appear error")
	QueryStringUnsupportedEndTime = errors.New("query string unsupported end time error")

	QueryStringMissingOffsetUnits     = errors.New("query string missing offset units error")
	QueryStringUnsupportedOffsetUnits = errors.New("query string unsupported offset units error")

	QueryStringMissingOffsetValue     = errors.New("query string missing offset value error")
	QueryStringUnsupportedOffsetValue = errors.New("query string unsupported offset value error")

	QueryStringMissingStartTime     = errors.New("query string missing start time error")
	QueryStringUnsupportedStartTime = errors.New("query string unsupported start time error")

	QueryStringTimeDifferenceTooBig = errors.New("query string difference between start_time and end_time is too big error")
	QueryStringEndTimePreviousThanStartTime = errors.New("query string end time is previous in time than start time error")

	QueryStringPnSn = errors.New("query string Product Number missing but Serial Number present error")
)
