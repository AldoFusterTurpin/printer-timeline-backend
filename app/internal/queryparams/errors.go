package queryparams

import (
	. "bitbucket.org/aldoft/printer-timeline-backend/app/internal/errorTypes"
)

const (
	ErrorQueryStringMissingTimeRangeType     = ConstError("query string missing time range type error")
	ErrorQueryStringUnsupportedTimeRangeType = ConstError("query string unsupported time range type error")

	ErrorQueryStringStartTimeAppears   = ConstError("query string start time should not appear error")
	ErrorQueryStringMissingEndTime     = ConstError("query string missing end time when time range is absolute error")
	ErrorQueryStringEndTimeAppears     = ConstError("query string end time should not appear error")
	ErrorQueryStringUnsupportedEndTime = ConstError("query string unsupported end time error")

	ErrorQueryStringMissingOffsetUnits     = ConstError("query string missing offset units error")
	ErrorQueryStringUnsupportedOffsetUnits = ConstError("query string unsupported offset units error")

	ErrorQueryStringMissingOffsetValue     = ConstError("query string missing offset value error")
	ErrorQueryStringUnsupportedOffsetValue = ConstError("query string unsupported offset value error")

	ErrorQueryStringMissingStartTime     = ConstError("query string missing start time error")
	ErrorQueryStringUnsupportedStartTime = ConstError("query string unsupported start time error")

	ErrorQueryStringTimeDifferenceTooBig         = ConstError("query string difference between start_time and end_time is too big error")
	ErrorQueryStringEndTimePreviousThanStartTime = ConstError("query string end time is previous in time than start time error")

	ErrorQueryStringPnSn = ConstError("query string Product Number missing but Serial Number present error")

	ErrorQueryStringMissingBucketRegion     = ConstError("query string missing bucket region error")
	ErrorQueryStringMissingBucketName       = ConstError("query string missing bucket name error")
	ErrorQueryStringMissingObjectKey        = ConstError("query string missing object key error")
	ErrorQueryStringUnsupportedBucketRegion = ConstError("query string unsupported bucket region error")

	ErrorNotValidEndpoint = ConstError("invalid endpoint reached")
)
