package queryParamsCtrl

import (
	"bitbucket.org/aldoft/printer-timeline-backend/errors"
	"strconv"
	"time"
)

func stringToTime(s string) (time.Time, error) {
	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0), nil
}

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
	return startTime, endTime, nil
}

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

func processAbsoluteTime(startTimeEpoch, endTimeEpoch string) (startTime time.Time, endTime time.Time, err error) {
	if startTimeEpoch == "" {
		return time.Time{}, time.Time{}, errors.QueryStringMissingStartTime
	}

	startTime, err = stringToTime(startTimeEpoch)
	if err != nil {
		return time.Time{}, time.Time{}, errors.QueryStringUnsupportedStartTime
	}

	if endTimeEpoch == "" {
		return time.Time{}, time.Time{}, errors.QueryStringMissingEndTime
	}
	endTime, err = stringToTime(endTimeEpoch)
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

func ExtractTimeRange(queryParameters map[string]string) (startTime time.Time, endTime time.Time, err error) {
	timeType := queryParameters["time_type"]
	if timeType == "" {
		return time.Time{}, time.Time{}, errors.QueryStringMissingTimeRangeType
	}

	startTimeEpoch := queryParameters["start_time"]
	endTimeEpoch :=  queryParameters["end_time"]
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
