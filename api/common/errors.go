package common

import "errors"

var (
	QueryStringMissingTimeRangeType     = errors.New("query string missing time range type error")
	QueryStringUnsupportedTimeRangeType = errors.New("query string unsupported time range type error")
	QueryStringMissingStartTime         = errors.New("query string missing start time when time range is absolute error")
	QueryStringMissingEndTime           = errors.New("query string missing end time when time range is absolute error")
	QueryStringPresentEndTime           = errors.New("query string end time should not appear error")
	QueryStringPnSn                     = errors.New("query string Product Number missing but Serial Number present")
)
