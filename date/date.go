package date

import (
	"time"
)

// dateFormat is the format which git expects for date ranges
const dateFormat = "2006-01-02"

// Yesterday returns yesterday's date
func Yesterday() string {
	return time.Now().AddDate(0, 0, -1).Format(dateFormat)
}

// Range returns a string string specifying date limits for git
func Range(floor string) ([]string, error) {
	d, err := time.Parse(dateFormat, floor)
	if err != nil {
		return nil, err
	}

	ceiling := d.AddDate(0, 0, 1).Format(dateFormat)

	dateRange := []string{
		"--since",
		floor + ":00:00",
		"--before",
		ceiling + ":00:00",
	}

	return dateRange, nil
}
