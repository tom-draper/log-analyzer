package display

import "fmt"

func DisplayLines(lines []map[string]any) {
	fmt.Printf("%d lines\n", len(lines))
	for i, line := range lines {
		fmt.Printf("line %d\n", i)
		for k, v := range line {
			fmt.Printf("	%s: %s\n", k, v)
		}
	}
}
