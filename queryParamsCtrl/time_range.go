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


func ExtractTimeRange(queryParameters map[string]string) (startTime time.Time, endTime time.Time, err error) {

	timeTypeStr := queryParameters["time_type"]
	if timeTypeStr == "" {
		err = errors.QueryStringMissingTimeRangeType
		return
	}

	startTimeEpoch := queryParameters["start_time"]
	endTimeEpoch := queryParameters["end_time"]
	offsetUnits := queryParameters["offset_units"]
	offsetValue := queryParameters["offset_value"]
	switch timeTypeStr {
	case "relative":
		if startTimeEpoch != "" {
			err = errors.QueryStringStartTimeAppears
			return
		}
		if endTimeEpoch != "" {
			err = errors.QueryStringEndTimeAppears
			return
		}
		if offsetUnits == "" {
			err = errors.QueryStringMissingOffsetUnits
			return
		}
		if offsetUnits != "seconds" && offsetUnits != "minutes" {
			err = errors.QueryStringUnsupportedOffsetUnits
			return
		}
		if offsetValue == "" {
			err = errors.QueryStringMissingOffsetValue
			return
		}

		var offsetValueInt int
		offsetValueInt, err = strconv.Atoi(offsetValue)
		if err != nil {
			err = errors.QueryStringUnsupportedOffsetValue
			return
		}

		if offsetUnits == "minutes" && offsetValueInt > 60 {
			err = errors.QueryStringUnsupportedOffsetValue
			return
		}
		if offsetUnits == "seconds" && offsetValueInt > 3600 {
			err = errors.QueryStringUnsupportedOffsetValue
			return
		}
		if offsetValueInt < 1 {
			err = errors.QueryStringUnsupportedOffsetValue
			return
		}

		endTime = time.Now()

		var durationOffset time.Duration
		if offsetUnits == "minutes" {
			durationOffset = -1 * time.Minute * time.Duration(offsetValueInt)
		} else if offsetUnits == "seconds" {
			durationOffset = -1 * time.Second * time.Duration(offsetValueInt)
		}
		startTime = time.Now().Add(durationOffset)

	case "absolute":
		if startTimeEpoch == "" {
			err = errors.QueryStringMissingStartTime
			return
		}

		startTime, err = stringToTime(startTimeEpoch)
		if err != nil {
			err = errors.QueryStringUnsupportedStartTime
			return
		}

		if endTimeEpoch == "" {
			err = errors.QueryStringMissingEndTime
			return
		}
		endTime, err = stringToTime(endTimeEpoch)
		if err != nil {
			err = errors.QueryStringUnsupportedEndTime
			return
		}

		diff := endTime.Sub(startTime)
		if diff.Minutes() > 60 {
			err = errors.QueryStringTimeDifferenceTooBig
			return
		}
		if diff.Minutes() < 0 {
			err = errors.QueryStringEndTimePreviousThanStartTime
			return
		}
	default:
		err = errors.QueryStringUnsupportedTimeRangeType
		return
	}
	return startTime, endTime, nil
}