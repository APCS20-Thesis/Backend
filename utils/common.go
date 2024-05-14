package utils

import "time"

func ToTimeString(t time.Time) string {
	defaultTime := time.Time{}
	if t == defaultTime {
		return ""
	}
	return t.String()
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
