package cronparser_test

import (
	"reflect"
	"testing"

	"github.com/disq/cronparser"
)

type testCase struct {
	Input string

	ExpectedValues map[string][]int // CronField Name vs. values
	ExpectedCmd    string
	ExpectedErr    error
}

func TestBasicCases(t *testing.T) {
	var cases = []testCase{
		{
			Input: "*/15 0 1,15 * Mon /usr/bin/find",

			ExpectedValues: map[string][]int{
				"minute":       {0, 15, 30, 45},
				"hour":         {0},
				"day of month": {1, 15},
				"month":        {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				"day of week":  {1},
			},
			ExpectedCmd: "/usr/bin/find",
			ExpectedErr: nil,
		},
		{
			Input: "*/15 0 1,15 * Mon-Fri /usr/bin/find",

			ExpectedValues: map[string][]int{
				"minute":       {0, 15, 30, 45},
				"hour":         {0},
				"day of month": {1, 15},
				"month":        {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				"day of week":  {1, 2, 3, 4, 5},
			},
			ExpectedCmd: "/usr/bin/find",
			ExpectedErr: nil,
		},
	}

	for _, cs := range cases {
		vals, cmd, err := cronparser.Parse(cs.Input)
		if err != cs.ExpectedErr {
			t.Errorf("Want err %v, got %v", cs.ExpectedErr, err)
		}
		if cmd != cs.ExpectedCmd {
			t.Errorf("Want cmd %v, got %v", cs.ExpectedCmd, cmd)
		}
		valMap := make(map[string][]int)
		for _, f := range vals {
			for _, v := range f.Values {
				valMap[f.Name] = append(valMap[f.Name], v)
			}
		}

		if !reflect.DeepEqual(valMap, cs.ExpectedValues) {
			t.Errorf("Want vals:\n%+v, got:\n %+v", cs.ExpectedValues, valMap)
		}
	}

}
