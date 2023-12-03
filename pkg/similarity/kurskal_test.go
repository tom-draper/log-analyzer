package similarity

import "testing"

var nodes = []Node{
	{0, "a"},
	{1, "b"},
	{2, "c"},
	{3, "d"},
	{4, "e"},
	{5, "f"},
	{6, "g"},
	{7, "h"},
	{8, "i"},
	{9, "j"},
	{10, "k"},
}

var edges = []Graph{
	{
		{1, &nodes[0], &nodes[1]},
		{2, &nodes[0], &nodes[2]},
		{3, &nodes[0], &nodes[3]},
		{4, &nodes[0], &nodes[4]},
		{5, &nodes[1], &nodes[2]},
		{6, &nodes[1], &nodes[3]},
		{7, &nodes[1], &nodes[4]},
		{8, &nodes[2], &nodes[3]},
		{9, &nodes[2], &nodes[4]},
		{10, &nodes[3], &nodes[4]},
		{11, &nodes[5], &nodes[6]},
		{12, &nodes[5], &nodes[7]},
		{13, &nodes[5], &nodes[8]},
		{14, &nodes[5], &nodes[9]},
		{15, &nodes[6], &nodes[7]},
		{16, &nodes[6], &nodes[8]},
		{17, &nodes[6], &nodes[9]},
		{18, &nodes[7], &nodes[8]},
		{19, &nodes[7], &nodes[9]},
		{20, &nodes[8], &nodes[9]},
		{21, &nodes[0], &nodes[5]},
		{22, &nodes[1], &nodes[6]},
		{23, &nodes[2], &nodes[7]},
		{24, &nodes[3], &nodes[8]},
		{25, &nodes[4], &nodes[9]},
		{26, &nodes[0], &nodes[6]},
		{27, &nodes[0], &nodes[7]},
		{28, &nodes[0], &nodes[8]},
		{29, &nodes[0], &nodes[9]},
		{30, &nodes[1], &nodes[7]},
		{31, &nodes[1], &nodes[8]},
		{32, &nodes[1], &nodes[9]},
		{33, &nodes[2], &nodes[8]},
		{34, &nodes[2], &nodes[9]},
		{35, &nodes[3], &nodes[9]},
		{36, &nodes[0], &nodes[10]},
		{37, &nodes[1], &nodes[10]},
		{38, &nodes[2], &nodes[10]},
		{39, &nodes[3], &nodes[10]},
		{40, &nodes[4], &nodes[10]},
		{41, &nodes[5], &nodes[10]},
	},
	{
		{99, &nodes[0], &nodes[1]},
		{63, &nodes[0], &nodes[2]},
		{53, &nodes[0], &nodes[3]},
		{20, &nodes[0], &nodes[4]},
		{25, &nodes[1], &nodes[2]},
		{19, &nodes[1], &nodes[3]},
		{14, &nodes[1], &nodes[4]},
		{57, &nodes[2], &nodes[3]},
		{22, &nodes[2], &nodes[4]},
	},
	{
		{1, &nodes[0], &nodes[1]},
		{0, &nodes[0], &nodes[2]},
		{0, &nodes[0], &nodes[3]},
		{0, &nodes[0], &nodes[4]},
		{0, &nodes[1], &nodes[2]},
		{0, &nodes[1], &nodes[3]},
		{0, &nodes[1], &nodes[4]},
		{0, &nodes[2], &nodes[3]},
		{0, &nodes[2], &nodes[4]},
		{0, &nodes[3], &nodes[4]},
		{1, &nodes[5], &nodes[6]},
		{2, &nodes[5], &nodes[7]},
		{3, &nodes[5], &nodes[8]},
		{4, &nodes[5], &nodes[9]},
		{5, &nodes[6], &nodes[7]},
		{6, &nodes[6], &nodes[8]},
		{7, &nodes[6], &nodes[9]},
		{8, &nodes[7], &nodes[8]},
		{9, &nodes[7], &nodes[9]},
	},
	{
		{0, &nodes[0], &nodes[1]},
		{0, &nodes[0], &nodes[2]},
		{0, &nodes[0], &nodes[3]},
		{0, &nodes[0], &nodes[4]},
		{0, &nodes[1], &nodes[2]},
		{0, &nodes[1], &nodes[3]},
		{0, &nodes[1], &nodes[4]},
		{0, &nodes[2], &nodes[3]},
		{0, &nodes[2], &nodes[4]},
		{0, &nodes[3], &nodes[4]},
		{0, &nodes[5], &nodes[6]},
		{0, &nodes[5], &nodes[7]},
		{0, &nodes[5], &nodes[8]},
		{0, &nodes[5], &nodes[9]},
		{0, &nodes[6], &nodes[7]},
		{0, &nodes[6], &nodes[8]},
		{0, &nodes[6], &nodes[9]},
		{0, &nodes[7], &nodes[8]},
		{0, &nodes[7], &nodes[9]},
	},
}

func numberOfNodes(graph Graph) int {
	nodes := make(map[int]struct{})
	for _, edge := range graph {
		nodes[edge.node1.id] = struct{}{}
		nodes[edge.node2.id] = struct{}{}
	}
	return len(nodes)
}

func TestKruskal(t *testing.T) {
	for i := range edges {
		mst := Kruskal(numberOfNodes(edges[i]), edges[i])

		for _, edge := range mst {
			t.Log(edge)
		}

		if len(mst) == 0 {
			t.Error("Kruskal failed to find minimum spanning tree")
		}
		if len(mst) != 5 {
			t.Errorf("Expected 5 edges, got %d", len(mst))
		}
	}
}
