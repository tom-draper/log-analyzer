package similarity

import "testing"

func TestKruskal(t *testing.T) {
	Kruskal(5, []Edge{
		{1, &Node{0, "a"}, &Node{1, "b"}},
		{2, &Node{0, "a"}, &Node{2, "c"}},
		{3, &Node{0, "a"}, &Node{3, "d"}},
		{4, &Node{0, "a"}, &Node{4, "e"}},
		{5, &Node{1, "b"}, &Node{2, "c"}},
		{6, &Node{1, "b"}, &Node{3, "d"}},
		{7, &Node{1, "b"}, &Node{4, "e"}},
		{8, &Node{2, "c"}, &Node{3, "d"}},
		{9, &Node{2, "c"}, &Node{4, "e"}},
		{10, &Node{3, "d"}, &Node{4, "e"}},
	})
}
