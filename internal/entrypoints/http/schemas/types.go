package schemas

import "time"

type Date string

// Returns string formatted as "dd-mm-yyyy".
func NewDate(t time.Time) Date {
	return Date(t.Format("02-01-2006"))
}

// ToTime converts the date string to time.Time.
func (d Date) ToTime() time.Time {
	layout := "02-01-2006"
	parsedTime, _ := time.Parse(layout, string(d))
	return parsedTime
}
