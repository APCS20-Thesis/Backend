package utils

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strings"
	"time"
)

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

func ConvertWrappersBoolToBoolAdd(val *wrapperspb.BoolValue) *bool {
	if val == nil {
		return nil
	}
	value := val.GetValue()
	return &value
}

func TransformPassword(password string) string {
	if len(password) < 6 {
		// Handle case where password is less than 6 characters
		return strings.Repeat("*", len(password))
	}
	return strings.Repeat("*", len(password)-6) + password[len(password)-6:]
}

func CompareDay(date1 time.Time, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
