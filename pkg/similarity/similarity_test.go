package similarity

import (
	"fmt"
	"testing"
)

func TestFindGroups(t *testing.T) {
	lines := []string{
		"Hello",
		"World",
		"hello",
		"world",
		"Hello World",
		"test",
		"testing",
		"test1",
		"Test",
	}
	groups := FindGroups(lines)
	fmt.Println(groups)
}
