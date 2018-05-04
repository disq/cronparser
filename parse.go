package cronparser

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	lineRe = regexp.MustCompile(`^\s*(.+?)\s+(.+?)\s+(.+?)\s+(.+?)\s+(.+?)\s+(.+?)\s*$`)
	starRe = regexp.MustCompile(`^(\*\/[0-9]+)$`)
	rangRe = regexp.MustCompile(`^([0-9]+)-([0-9]+)$`)
)

func Parse(s string) (values []CronValue, cmd string, retErr error) {
	fields := Fields()

	ma := lineRe.FindStringSubmatch(s)
	if len(ma) != len(fields)+2 { // ma[0] is complete string, ma[last] is command...
		retErr = fmt.Errorf("Format mismatch")
		return
	}
	cmd = ma[len(fields)+1]

	for i := 0; i < len(fields); i++ {
		v := ma[i+1]
		possibleVals, err := fields[i].parsePossibleValues(v)
		if err != nil {
			retErr = fmt.Errorf("Parsing %v: %q is invalid: %v", fields[i].Name, v, err)
			return
		}

		values = append(values, CronValue{
			CronField: fields[i],
			Values:    possibleVals,
		})
	}

	return
}

func (f *CronField) parsePossibleValues(s string) (vals []int, retErr error) {
	defer func() {
		if retErr == nil && len(vals) == 0 {
			retErr = fmt.Errorf("No possible values")
		}
	}()

	commaParts := strings.Split(s, ",")
	for _, part := range commaParts {
		if ma := starRe.FindStringSubmatch(part); len(ma) == 2 {
			v, err := f.parseStar(ma[1])
			if err != nil {
				retErr = err
				return
			}
			vals = append(vals, v...)
			continue
		}

		if ma := rangRe.FindStringSubmatch(part); len(ma) == 3 {
			v, err := f.parseRange(ma[1], ma[2])
			if err != nil {
				retErr = err
				return
			}
			vals = append(vals, v...)
			continue
		}

	}

	return
}

func (f *CronField) parseStar(i string) (vals []int, retErr error) {
	return
}

func (f *CronField) parseRange(start, end string) (vals []int, retErr error) {
	return
}
