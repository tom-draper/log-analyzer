package parse

import (
	"fmt"
	"reflect"
)

func DisplayLines(extraction []Extraction) {
	fmt.Printf("%d lines\n", len(extraction))
	for i, e := range extraction {
		fmt.Printf("line %d\n", i+1)
		for k, v := range e.Params {
			fmt.Printf("	%s(%s): %v\n", k, reflect.TypeOf(v).String(), v)
		}
	}
}

func DisplayTestLines(extraction []Extraction) {
	fmt.Printf("%d lines\n", len(extraction))
	for _, e := range extraction {
		fmt.Printf("line %d: %s\n", e.LineNumber+1, e.Line)
		for k, v := range e.Params {
			fmt.Printf("	%s(%s): %v\n", k, reflect.TypeOf(v).String(), v)
		}
	}
}
