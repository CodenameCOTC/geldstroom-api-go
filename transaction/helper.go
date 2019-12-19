package transaction

import (
	"time"
)

type Range struct {
	firstDay time.Time
	lastDay  time.Time
}

var layout = "2006/01/02"

func getRange(date, dateRange string) (*Range, error) {
	if dateRange == oneDayRange {
		return getOneDayRange(date)
	} else if dateRange == oneWeekRange {
		return getOneWeekRange(date)
	} else if dateRange == oneMonthRange {
		return getOneMonthRange(date)
	} else {
		return nil, errInvalidDateRange
	}
}

func getOneDayRange(date string) (*Range, error) {
	var r Range

	t, err := time.Parse(layout, date)
	if err != nil {
		return nil, errInvalidDateFormat
	}

	r.firstDay = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	r.lastDay = r.firstDay.AddDate(0, 0, 1).Add(time.Nanosecond * -1)

	return &r, nil

}

func getOneWeekRange(date string) (*Range, error) {
	var r Range

	t, err := time.Parse(layout, date)
	if err != nil {
		return nil, errInvalidDateFormat
	}

	r.firstDay = time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday()), 0, 0, 0, 0, time.UTC)
	r.lastDay = r.firstDay.AddDate(0, 0, 7).Add(time.Nanosecond * -1)

	return &r, nil
}

func getOneMonthRange(date string) (*Range, error) {
	var r Range

	t, err := time.Parse(layout, date)
	if err != nil {
		return nil, errInvalidDateFormat
	}

	r.firstDay = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	r.lastDay = r.firstDay.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

	return &r, nil
}
