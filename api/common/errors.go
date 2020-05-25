package common

import "errors"

var (
	QueryStringMissingTimeRangeTypeError     = errors.New("query string missing time range type error")
	QueryStringUnsupportedTimeRangeTypeError = errors.New("query string unsupported time range type error")
	QueryStringStartTimeAppearsError         = errors.New("query string start time should not appear error")
	QueryStringEndTimeAppearsError           = errors.New("query string end time should not appear error")
	QueryStringMissingOffsetUnitsError       = errors.New("query string missing offset units error")
	QueryStringMissingOffsetValueError       = errors.New("query string missing offset value error")
	QueryStringUnsupportedOffsetUnitsError   = errors.New("query string unsupported offset units error")
	QueryStringUnsupportedOffsetValueError   = errors.New("query string unsupported offset value error")

	QueryStringMissingStartTimeError         = errors.New("query string missing start time") // when time range is absolute error")
	QueryStringMissingEndTimeError           = errors.New("query string missing end time when time range is absolute error")
	QueryStringPnSnError                     = errors.New("query string Product Number missing but Serial Number present")
)
