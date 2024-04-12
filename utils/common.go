package utils

import "time"

func ToTimeString(t time.Time) string {
	defaultTime := time.Time{}
	if t == defaultTime {
		return ""
	}
	return t.String()
}
