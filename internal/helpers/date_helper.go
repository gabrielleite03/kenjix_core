package helpers

import "time"

func ParseDate(value string) time.Time {
	t, _ := time.Parse("2006-01-02", value)
	return t
}
