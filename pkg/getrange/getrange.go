package getrange

import (
	"fmt"
	"time"
)

type Range struct {
	FirstDay time.Time
	LastDay  time.Time
}

var (
	Layout               = "2006/01/02"
	OneDayRange          = "DAILY"
	OneWeekRange         = "WEEKLY"
	OneMonthRange        = "MONTHLY"
	ErrInvalidDateRange  = fmt.Errorf("Date range must be one of %v | %v | %v", OneDayRange, OneWeekRange, OneMonthRange)
	ErrInvalidDateFormat = fmt.Errorf("Date format must be : %v", Layout)
)

func GetRange(date, dateRange string) (*Range, error) {
	if dateRange == OneDayRange {
		return GetOneDayRange(date)
	} else if dateRange == OneWeekRange {
		return GetOneWeekRange(date)
	} else if dateRange == OneMonthRange {
		return GetOneMonthRange(date)
	} else {
		return nil, ErrInvalidDateRange
	}
}

func GetOneDayRange(date string) (*Range, error) {
	var r Range

	t, err := time.Parse(Layout, date)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	r.FirstDay = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	r.LastDay = r.FirstDay.AddDate(0, 0, 1).Add(time.Nanosecond * -1)

	return &r, nil

}

func GetOneWeekRange(date string) (*Range, error) {
	var r Range

	t, err := time.Parse(Layout, date)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	r.FirstDay = time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday()), 0, 0, 0, 0, time.UTC)
	r.LastDay = r.FirstDay.AddDate(0, 0, 7).Add(time.Nanosecond * -1)

	return &r, nil
}

func GetOneMonthRange(date string) (*Range, error) {
	var r Range

	t, err := time.Parse(Layout, date)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	r.FirstDay = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	r.LastDay = r.FirstDay.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

	return &r, nil
}
