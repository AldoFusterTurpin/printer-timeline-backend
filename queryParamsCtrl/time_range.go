package queryParamsCtrl

import (
	"bitbucket.org/aldoft/printer-timeline-backend/errors"
	"strconv"
	"time"
)

func ExtractTimeRange(queryParameters map[string]string) (startTime int64, endTime int64, err error) {

	timeTypeStr := queryParameters["time_type"]
	if timeTypeStr == "" {
		err = errors.QueryStringMissingTimeRangeType
		return
	}

	startTimeStr := queryParameters["start_time"]
	endTimeStr := queryParameters["end_time"]
	offsetUnits := queryParameters["offset_units"]
	offsetValue := queryParameters["offset_value"]
	switch timeTypeStr {
	case "relative":
		if startTimeStr != "" {
			err = errors.QueryStringStartTimeAppears
			return
		}
		if endTimeStr != "" {
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

		endTime = time.Now().Unix()

		var duration time.Duration
		if offsetUnits == "minutes" {
			duration = -1 * time.Minute * time.Duration(offsetValueInt)
		} else if offsetUnits == "seconds" {
			duration = -1 * time.Second * time.Duration(offsetValueInt)
		}
		startTime = time.Now().Add(duration).Unix()

	case "absolute":
		if startTimeStr == "" {
			err = errors.QueryStringMissingStartTime
			return
		}

		startTime, err = strconv.ParseInt(startTimeStr, 10, 64)
		if err != nil {
			err = errors.QueryStringUnsupportedStartTime
			return
		}

		if endTimeStr == "" {
			err = errors.QueryStringMissingEndTime
			return
		}
		endTime, err = strconv.ParseInt(endTimeStr, 10, 64)
		if err != nil {
			err = errors.QueryStringUnsupportedEndTime
			return
		}

		diff := time.Unix(endTime, 0).Sub(time.Unix(startTime, 0))
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