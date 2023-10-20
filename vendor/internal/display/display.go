package display

import (
	"fmt"
	"reflect"
)

// DisplayLines prints all extracted parameters from log text.
func DisplayLines(params []map[string]any) {
	fmt.Printf("%d lines\n", len(params))
	for i, p := range params {
		fmt.Printf("line %d\n", i)
		for k, v := range p {
			fmt.Printf("	%s(%s): %v\n", k, reflect.TypeOf(v).String(), v)
		}
	}
}

// DisplayTestLines prints a sample of extracted parameters along with the
// origin log text line and line number to allow for user evaluation.
func DisplayTestLines(lines []string, params []map[string]any, indicies []int) {
	fmt.Printf("%d lines\n", len(lines))
	for i, p := range params {
		fmt.Printf("line %d: %s\n", indicies[i], lines[i])
		for k, v := range p {
			fmt.Printf("	%s(%s): %v\n", k, reflect.TypeOf(v).String(), v)
		}
	}
}
