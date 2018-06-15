package cronparser

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	lineRe = regexp.MustCompile(`^\s*(.+?)\s+(.+?)\s+(.+?)\s+(.+?)\s+(.+?)\s+(.+?)\s*$`)
	rangRe = regexp.MustCompile(`^((.+?)-(.+?)|\*)(\/([0-9]+))?$`)
)

// Parse a crontab line. Returns []CronValue, the command, and error.
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

// parsePossibleValues splits the non-whitespace string s into comma-separated parts and attempts to parse each part
func (f *CronField) parsePossibleValues(s string) (vals []int, retErr error) {
	defer func() {
		if retErr == nil && len(vals) == 0 {
			retErr = fmt.Errorf("No possible values")
		}
	}()

	commaParts := strings.Split(s, ",")
	for _, part := range commaParts {
		if ma := rangRe.FindStringSubmatch(part); len(ma) == 4 || len(ma) == 6 {
			if ma[1] == "*" { // "*" means first-last
				ma[2] = strconv.Itoa(f.Min)
				ma[3] = strconv.Itoa(f.Max)
			}
			step := ""
			if len(ma) == 6 {
				step = ma[5]
			}

			v, err := f.parseRange(ma[2], ma[3], step)
			if err != nil {
				retErr = err
				return
			}
			vals = append(vals, v...)
			continue
		}

		// Single digit without range
		result, err := f.parseNumeric(part)
		if err != nil {
			retErr = err
			return
		}

		vals = append(vals, result)
	}

	return
}

func (f *CronField) parseNumeric(s string) (result int, retErr error) {
	defer func() {
		if retErr != nil && f.Parser != nil { // About to return error? Try the custom parser...
			v := f.Parser(f, s)
			if v != nil { // Override error if indeed valid
				result = *v
				retErr = nil
			}
		}
	}()

	// Must be int
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		retErr = fmt.Errorf("not an integer")
		return
	}
	result = int(i) // int64 to int
	if result < f.Min || result > f.Max {
		retErr = fmt.Errorf("%v: out of range", result)
	}
	return
}

func (f *CronField) parseRange(start, end, step string) (vals []int, retErr error) {
	stp := 1
	if step != "" {
		if i, err := strconv.ParseInt(step, 10, 64); err != nil {
			retErr = fmt.Errorf("step should be numeric")
			return
		} else if i < 1 {
			retErr = fmt.Errorf("step should be positive")
			return
		} else {
			stp = int(i)
		}
	}

	var st, en int

	st, retErr = f.parseNumeric(start)
	if retErr != nil {
		return
	}

	en, retErr = f.parseNumeric(end)
	if retErr != nil {
		return
	}

	// Edge case: 23-4: 23..[f.Max] + [f.Min]..4
	if en < st {
		if stp != 1 {
			return nil, fmt.Errorf("Custom step not supported with reverse bounds")
		}
		for i := st; i <= f.Max; i += stp {
			vals = append(vals, i)
		}
		// TODO handle carry of stp
		for i := f.Min; i <= en; i += stp {
			vals = append(vals, i)
		}

		sort.Ints(vals)

		return
	}

	for i := st; i <= en; i += stp {
		vals = append(vals, i)
	}

	return
}
