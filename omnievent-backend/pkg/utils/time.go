package utils

import (
	"strconv"
	"time"
)

func GetUnixTime() int64 {
	return time.Now().Unix()
}

func GetUnixTimeFromTime(t time.Time) int64 {
	return t.Unix()
}

func GetMillisTime() int64 {
	return time.Now().UnixMilli()
}

func GetCurrentTime() time.Time {
	return time.Now()
}

func FromUnixTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

func AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

func FormatUnixTime(unix int64, layout string) string {
	if unix <= 0 {
		return ""
	}
	return FromUnixTime(unix).Format(layout)
}

func ParseTime(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func StringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func StringToIntWithDefault(s string, defaultVal int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return i
}

func StringToInt64WithDefault(s string, defaultVal int64) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultVal
	}
	return i
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
