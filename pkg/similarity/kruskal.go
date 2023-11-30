package similarity

import (
	"fmt"
	"sort"
)

func findParent(parent []int, i int) int {
	if parent[i] == i {
		return parent[i]
	}
	return findParent(parent, parent[i])
}

func Kruskal(v int, edges []Edge) []Edge {
	result := make([]Edge, 0)

	// Sort list of connections by their weight field
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].weight < edges[j].weight
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
		u, v, weight := edges[i].node1, edges[i].node2, edges[i].weight
		i += 1
		x := findParent(parent, u.id)
		y := findParent(parent, v.id)

		if x != y {
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
	}

	for _, edge := range result {
		fmt.Println(edge.String())
	}

	return result
}
