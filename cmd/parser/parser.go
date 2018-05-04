package main

import (
	"fmt"
	"os"
	"strings"

	"strconv"

	"github.com/disq/cronparser"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %v [cron line...]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	var input string

	if l := len(os.Args); l == 1 {
		printUsage()
	}

	if len(os.Args) == 2 {
		input = os.Args[1]
	} else {
		input = strings.Join(os.Args[1:], " ")
	}

	if strings.HasPrefix(strings.TrimLeft(input, "-"), "h") { // -help, --help, ---help, -h, --h, ...
		printUsage()
	}

	vals, cmd, err := cronparser.Parse(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	const format = "%14v %v\n"

	for _, v := range vals {
		s := make([]string, len(v.Values))
		for i, j := range v.Values {
			s[i] = strconv.Itoa(j)
		}
		fmt.Printf(format, v.Name, strings.Join(s, " "))
	}
	fmt.Printf(format, "command", cmd)
}
