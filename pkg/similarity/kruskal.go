package similarity

import (
	"sort"
)

func findParent(parent []int, i int) int {
	if parent[i] == i {
		return parent[i]
	}
	return findParent(parent, parent[i])
}

func Kruskal(v int, graph Graph) Graph {
	result := make(Graph, 0)

	// Sort list of connections by their weight field
	sort.Slice(graph, func(i, j int) bool {
		return graph[i].weight < graph[j].weight
	})

	parent := make([]int, v)
	rank := make([]int, v)
	// loop v times
	for i := 0; i < v; i++ {
		parent[i] = i
		rank[i] = 0
	}

	i := 0
	e := 0
	for e < v-1 {
		u, v, weight := graph[i].node1, graph[i].node2, graph[i].weight
		i += 1
		x := findParent(parent, u.id)
		y := findParent(parent, v.id)

		if x == y {
			continue
		}

		e += 1
		result = append(result, Edge{weight, u, v})
		// Union of the two sets x and y
		if rank[x] > rank[y] {
			parent[y] = x
		} else if rank[x] < rank[y] {
			parent[x] = y
		} else {
			parent[x] = y
			rank[y] += 1
		}
	}

	return result
}
