package cronparser

import (
	"strings"
	"time"
)

// CronField represents a single field in a crontab line
type CronField struct {
	Name   string
	Min    int
	Max    int
	Parser fieldParserFunc
}

// CronValue is possible values per CronField
type CronValue struct {
	CronField
	Values []int
}

type fieldParserFunc func(*CronField, string) *int

// Fields returns valid CronFields in a standard crontab line
func Fields() []CronField {
	return []CronField{
		{"minute", 0, 59, nil},
		{"hour", 0, 23, nil},
		{"day of month", 1, 31, nil},
		{"month", 1, 12, parseMonth},
		{"day of week", 0, 6, parseWeekday},
	}
}

// parseMonth parses a 3-letter month to index, or nil.
func parseMonth(f *CronField, s string) *int {
	s = strings.ToLower(s)

	for i := f.Min; i <= f.Max; i++ {
		if s == strings.ToLower(time.Month(i).String()[:3]) {
			return &i
		}
	}

	return nil
}

// parseWeekday parses a 3-letter weekday to index, or nil.
func parseWeekday(f *CronField, s string) *int {
	if s == "7" { // Sun is both 0 and 7
		i := 0
		return &i
	}

	s = strings.ToLower(s)

	for i := f.Min; i <= f.Max; i++ {
		if s == strings.ToLower(time.Weekday(i).String()[:3]) {
			return &i
		}
	}

	return nil
}
