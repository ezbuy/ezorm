package db

import "time"

func TimeToLocalTime(c time.Time) string {
	return c.Local().Format("2006-01-02 15:04:05")
}

func TimeParseLocalTime(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return t
	}
	localTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(),
		t.Second(), t.Nanosecond(), time.Local)
	return localTime
}
