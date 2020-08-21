// package errors defines the different errors of the application.
package errors

type constError string

func (err constError) Error() string {
	return string(err)
}

const (
	QueryStringMissingTimeRangeType     = constError("query string missing time range type error")
	QueryStringUnsupportedTimeRangeType = constError("query string unsupported time range type error")

	QueryStringStartTimeAppears   = constError("query string start time should not appear error")
	QueryStringMissingEndTime     = constError("query string missing end time when time range is absolute error")
	QueryStringEndTimeAppears     = constError("query string end time should not appear error")
	QueryStringUnsupportedEndTime = constError("query string unsupported end time error")

	QueryStringMissingOffsetUnits     = constError("query string missing offset units error")
	QueryStringUnsupportedOffsetUnits = constError("query string unsupported offset units error")

	QueryStringMissingOffsetValue     = constError("query string missing offset value error")
	QueryStringUnsupportedOffsetValue = constError("query string unsupported offset value error")

	QueryStringMissingStartTime     = constError("query string missing start time error")
	QueryStringUnsupportedStartTime = constError("query string unsupported start time error")

	QueryStringTimeDifferenceTooBig         = constError("query string difference between start_time and end_time is too big error")
	QueryStringEndTimePreviousThanStartTime = constError("query string end time is previous in time than start time error")

	QueryStringPnSn = constError("query string Product Number missing but Serial Number present error")

	QueryStringMissingBucketRegion     = constError("query string missing bucket region error")
	QueryStringMissingBucketName       = constError("query string missing bucket name error")
	QueryStringMissingObjectKey        = constError("query string missing object key error")
	QueryStringUnsupportedBucketRegion = constError("query string unsupported bucket region error")
)
