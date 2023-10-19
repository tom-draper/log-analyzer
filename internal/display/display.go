package display

import (
	"fmt"
	"reflect"
)

func DisplayLines(lines []map[string]any) {
	fmt.Printf("%d lines\n", len(lines))
	for i, line := range lines {
		fmt.Printf("line %d\n", i)
		for k, v := range line {
			fmt.Printf("	%s(%s): %v\n", k, reflect.TypeOf(v).String(), v)
		}
	}
}

func DisplayTestLines(lines []string, params []map[string]any, indicies []int) {
	fmt.Printf("%d lines\n", len(lines))
	for i, p := range params {
		fmt.Printf("line %d: %s\n", indicies[i], lines[i])
		for k, v := range p {
			fmt.Printf("	%s(%s): %v\n", k, reflect.TypeOf(v).String(), v)
		}
	}
}
