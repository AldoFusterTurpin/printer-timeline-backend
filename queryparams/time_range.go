package queryparams

import (
	"strconv"
	"time"

	"bitbucket.org/aldoft/printer-timeline-backend/errors"
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

// processOffset returns and start time and end time based in offsetUnits and offsetValue.
// It is used for the relative time in the api.
// It also returns an error if any.
func processOffset(offsetUnits, offsetValue string) (startTime time.Time, endTime time.Time, err error) {
	if offsetValue == "" {
		return time.Time{}, time.Time{}, errors.QueryStringMissingOffsetValue
	}

	offsetValueInt, err := strconv.Atoi(offsetValue)
	if err != nil {
		return time.Time{}, time.Time{}, errors.QueryStringUnsupportedOffsetValue
	}

	var duration time.Duration
	if offsetUnits == "minutes" {
		if offsetValueInt > 60 {
			return time.Time{}, time.Time{}, errors.QueryStringUnsupportedOffsetValue
		}
		duration = time.Minute

	} else if offsetUnits == "seconds" {
		if offsetValueInt > 3600 {
			return time.Time{}, time.Time{}, errors.QueryStringUnsupportedOffsetValue
		}
		duration = time.Second
	}

	if offsetValueInt < 1 {
		return time.Time{}, time.Time{}, errors.QueryStringUnsupportedOffsetValue
	}

	endTime = time.Now()
	durationOffset := duration * time.Duration(offsetValueInt)
	startTime = endTime.Add(-durationOffset)
	return startTime.UTC(), endTime.UTC(), nil
}

// processRelativeTime receives and start and end time as strings, ofsset untis and value and returns the appropiate start and end time.
// It also returns an error if any.
func processRelativeTime(startTimeEpoch, endTimeEpoch, offsetUnits, offsetValue string) (startTime time.Time, endTime time.Time, err error) {
	if startTimeEpoch != "" {
		return time.Time{}, time.Time{}, errors.QueryStringStartTimeAppears
	}

	if endTimeEpoch != "" {
		return time.Time{}, time.Time{}, errors.QueryStringEndTimeAppears
	}

	if offsetUnits == "" {
		return time.Time{}, time.Time{}, errors.QueryStringMissingOffsetUnits
	}

	if offsetUnits != "minutes" && offsetUnits != "seconds" {
		return time.Time{}, time.Time{}, errors.QueryStringUnsupportedOffsetUnits
	}
	return processOffset(offsetUnits, offsetValue)
}

// processRelativeTime receives and start and end timeas strings and returns the appropiate start and end time variables.
// It also returns an error if any.
func processAbsoluteTime(startTimeEpoch, endTimeEpoch string) (startTime time.Time, endTime time.Time, err error) {
	if startTimeEpoch == "" {
		return time.Time{}, time.Time{}, errors.QueryStringMissingStartTime
	}

	startTime, err = stringEpochToUTCTime(startTimeEpoch)
	if err != nil {
		return time.Time{}, time.Time{}, errors.QueryStringUnsupportedStartTime
	}

	if endTimeEpoch == "" {
		return time.Time{}, time.Time{}, errors.QueryStringMissingEndTime
	}
	endTime, err = stringEpochToUTCTime(endTimeEpoch)
	if err != nil {
		return time.Time{}, time.Time{}, errors.QueryStringUnsupportedEndTime
	}

	diff := endTime.Sub(startTime)
	if diff.Minutes() > 60 {
		return time.Time{}, time.Time{}, errors.QueryStringTimeDifferenceTooBig
	}
	if diff.Minutes() < 0 {
		return time.Time{}, time.Time{}, errors.QueryStringEndTimePreviousThanStartTime
	}

	return startTime, endTime, nil
}

// ExtractTimeRange extracts from the query parameters the appropiate start time and end time based in some
// logic using start_time, end_time, offset_units and offset_value.
// It also returns an error if any.
func ExtractTimeRange(queryParameters map[string]string) (startTime time.Time, endTime time.Time, err error) {
	timeType := queryParameters["time_type"]
	if timeType == "" {
		return time.Time{}, time.Time{}, errors.QueryStringMissingTimeRangeType
	}

	startTimeEpoch := queryParameters["start_time"]
	endTimeEpoch := queryParameters["end_time"]
	switch timeType {
	case "relative":
		offsetUnits := queryParameters["offset_units"]
		offsetValue := queryParameters["offset_value"]

		startTime, endTime, err = processRelativeTime(startTimeEpoch, endTimeEpoch, offsetUnits, offsetValue)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return startTime, endTime, nil

	case "absolute":
		startTime, endTime, err = processAbsoluteTime(startTimeEpoch, endTimeEpoch)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return startTime, endTime, nil

	default:
		return time.Time{}, time.Time{}, errors.QueryStringUnsupportedTimeRangeType
	}
}
