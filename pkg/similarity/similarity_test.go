package similarity

import (
	"fmt"
	"testing"
)

var demoLines = [][]string{
	{
		"Hello",
		"World",
		"hello",
		"world",
		"Hello World",
		"test",
		"testing",
		"test1",
		"Test",
	},
	{
		"a",
		"b",
		"aa",
		"bbb",
		"aba",
		"aaa",
		"c",
		"cc",
		"abc",
	},
	{
		"test",
		"test",
		"test",
		"test",
		"test",
		"testing",
		"Testing",
		"Test",
		"Test",
		"Test",
		"Test",
		"Test",
	},
	{
		"cat dog mouse elephant deer bear",
		"cat dog mouse elephant",
		"deer bear",
		"car bear",
		"car bear",
		"car bear elephant mouse",
		"car bear elephant deer",
		"dog",
		"dog cat",
		"beat deet elephant moust dot cat",
		"bear deer elephant mouse dog cat",
		"bear deer elephant mouse dog",
		"elephant",
		"elehant",
		"god tac",
		"god",
		"god dog",
		"goddog",
		"goddogcat",
		"goddogcatelepahnt",
		"bear deer",
		"bedeerar",
		"animal",
	},
	{
		"1",
		"2",
		"122",
		"123",
		"3",
		"4",
		"43",
		"34",
		"33",
		"333",
		"3333",
		"124",
		"244",
		"2444",
		"135",
		"321",
		"12345",
	},
}

func TestFindGroups(t *testing.T) {
	for _, lines := range demoLines {
		groups := FindGroupsFromLines(lines)
		fmt.Println(groups)
	}
}
