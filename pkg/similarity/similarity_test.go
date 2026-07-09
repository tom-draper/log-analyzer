package similarity

import (
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

func TestFindGroupsSeparatesDistinctFormats(t *testing.T) {
	lines := []string{
		`[11/Dec/2023:11:01:28] 220.203.23.174 "GET /home HTTP/1.1" 200 182`,
		`[11/Dec/2023:11:01:29] 89.238.65.53 "POST /users HTTP/1.1" 201 145`,
		`[11/Dec/2023:11:01:30] 209.51.141.74 "GET /about HTTP/1.1" 404 320`,
		`ERROR: database connection timed out after 30s`,
		`ERROR: database connection refused by host db-2`,
		`ERROR: database connection reset while querying users`,
	}

	groups := FindGroupsFromLines(lines)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups (access log vs error), got %d: %+v", len(groups), groups)
	}

	// Groups are returned largest first; here both have 3 members.
	for _, g := range groups {
		if len(g.Members) != 3 {
			t.Errorf("expected 3 members per group, got %d: %v", len(g.Members), g.Members)
		}
		// Each group should be internally consistent: all access-log lines
		// start with "[", all error lines start with "ERROR".
		prefix := g.Members[0][:1]
		for _, m := range g.Members {
			if m[:1] != prefix {
				t.Errorf("group mixes formats: %v", g.Members)
				break
			}
		}
	}
}

func TestFindGroups(t *testing.T) {
	for tc, lines := range demoLines {
		valid := make(map[string]struct{}, len(lines))
		for _, line := range lines {
			valid[line] = struct{}{}
		}

		groups := FindGroupsFromLines(lines)
		if len(groups) == 0 {
			t.Errorf("case %d: expected at least one group, got none", tc)
		}

		for g, group := range groups {
			if len(group.Members) == 0 {
				t.Errorf("case %d group %d: group has no members", tc, g)
			}
			// Every grouped member must be one of the original input lines.
			for _, member := range group.Members {
				if _, ok := valid[member]; !ok {
					t.Errorf("case %d group %d: member %q is not an input line", tc, g, member)
				}
			}
		}
	}
}
