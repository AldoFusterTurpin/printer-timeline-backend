package queryparams

import (
	"strconv"
	"time"
)

// stringEpochToUTCTime converts an epoch string to the corresponding time in UTC.
func stringEpochToUTCTime(s string) (time.Time, error) {
	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	t := time.Unix(sec, 0)
	loc, _ := time.LoadLocation("UTC")
	return t.In(loc), nil
}

// processOffset returns an start time and end time based in offsetUnits and offsetValue.
// It is used for the relative time in the api.
// It also returns an error if any.
func processOffset(offsetUnits, offsetValue string, maxTimeDiffInMinutes int) (startTime time.Time, endTime time.Time, err error) {
	if offsetValue == "" {
		return time.Time{}, time.Time{}, ErrorQueryStringMissingOffsetValue
	}

	offsetValueInt, err := strconv.Atoi(offsetValue)
	if err != nil {
		return time.Time{}, time.Time{}, ErrorQueryStringUnsupportedOffsetValue
	}

	var duration time.Duration
	if offsetUnits == "minutes" {
		if offsetValueInt > maxTimeDiffInMinutes {
			return time.Time{}, time.Time{}, ErrorQueryStringUnsupportedOffsetValue
		}
		duration = time.Minute

	} else if offsetUnits == "seconds" {
		if offsetValueInt > 3600 {
			return time.Time{}, time.Time{}, ErrorQueryStringUnsupportedOffsetValue
		}
		duration = time.Second
	}

	if offsetValueInt < 1 {
		return time.Time{}, time.Time{}, ErrorQueryStringUnsupportedOffsetValue
	}

	endTime = time.Now()
	durationOffset := duration * time.Duration(offsetValueInt)
	startTime = endTime.Add(-durationOffset)
	return startTime.UTC(), endTime.UTC(), nil
}

// processRelativeTime receives and start and end time as strings, ofsset untis and value and returns the appropiate start and end time.
// It also returns an error if any.
func processRelativeTime(startTimeEpoch, endTimeEpoch, offsetUnits, offsetValue string, maxTimeDiffInMinutes int) (startTime time.Time, endTime time.Time, err error) {
	if startTimeEpoch != "" {
		return time.Time{}, time.Time{}, ErrorQueryStringStartTimeAppears
	}

	if endTimeEpoch != "" {
		return time.Time{}, time.Time{}, ErrorQueryStringEndTimeAppears
	}

	if offsetUnits == "" {
		return time.Time{}, time.Time{}, ErrorQueryStringMissingOffsetUnits
	}

	if offsetUnits != "minutes" && offsetUnits != "seconds" {
		return time.Time{}, time.Time{}, ErrorQueryStringUnsupportedOffsetUnits
	}
	return processOffset(offsetUnits, offsetValue, maxTimeDiffInMinutes)
}

// processRelativeTime receives and start and end timeas strings and returns the appropiate start and end time variables.
// It also returns an error if any.
func processAbsoluteTime(startTimeEpoch, endTimeEpoch string, maxTimeDiffInMinutes int) (startTime time.Time, endTime time.Time, err error) {
	if startTimeEpoch == "" {
		return time.Time{}, time.Time{}, ErrorQueryStringMissingStartTime
	}

	startTime, err = stringEpochToUTCTime(startTimeEpoch)
	if err != nil {
		return time.Time{}, time.Time{}, ErrorQueryStringUnsupportedStartTime
	}

	if endTimeEpoch == "" {
		return time.Time{}, time.Time{}, ErrorQueryStringMissingEndTime
	}
	endTime, err = stringEpochToUTCTime(endTimeEpoch)
	if err != nil {
		return time.Time{}, time.Time{}, ErrorQueryStringUnsupportedEndTime
	}

	diff := endTime.Sub(startTime)
	if diff.Minutes() > float64(maxTimeDiffInMinutes) {
		return time.Time{}, time.Time{}, ErrorQueryStringTimeDifferenceTooBig
	}
	if diff.Minutes() < 0 {
		return time.Time{}, time.Time{}, ErrorQueryStringEndTimePreviousThanStartTime
	}

	return startTime, endTime, nil
}

// ExtractTimeRange extracts from the query parameters the appropiate start time and end time based in some
// logic using start_time, end_time, offset_units and offset_value.
// It also returns an error if any.
func ExtractTimeRange(queryParameters map[string]string, maxTimeDiffInMinutes int) (startTime time.Time, endTime time.Time, err error) {
	timeType := queryParameters["time_type"]
	if timeType == "" {
		return time.Time{}, time.Time{}, ErrorQueryStringMissingTimeRangeType
	}

	startTimeEpoch := queryParameters["start_time"]
	endTimeEpoch := queryParameters["end_time"]
	switch timeType {
	case "relative":
		offsetUnits := queryParameters["offset_units"]
		offsetValue := queryParameters["offset_value"]

		startTime, endTime, err = processRelativeTime(startTimeEpoch, endTimeEpoch, offsetUnits, offsetValue, maxTimeDiffInMinutes)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return startTime, endTime, nil

	case "absolute":
		startTime, endTime, err = processAbsoluteTime(startTimeEpoch, endTimeEpoch, maxTimeDiffInMinutes)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return startTime, endTime, nil

	default:
		return time.Time{}, time.Time{}, ErrorQueryStringUnsupportedTimeRangeType
	}
}
