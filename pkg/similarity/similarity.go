package similarity

import (
	"fmt"

	"github.com/adrg/strutil/metrics"
)

type Group struct {
	Members []string
}

type Node struct {
	id    int
	value string
}

func (n Node) String() string {
	return fmt.Sprintf("[%d]%s", n.id, n.value)
}

type Edge struct {
	weight int
	node1  *Node
	node2  *Node
}

func (e Edge) String() string {
	return fmt.Sprintf("[%d]%s [%d]%s = %d", e.node1.id, e.node1.value, e.node2.id, e.node2.value, e.weight)
}

func extractGroups(edges []Edge) []Group {
	return []Group{}
}

func FindGroups(lines []string) []Group {
	// Create node for each line
	nodes := make([]Node, len(lines))
	for i, line := range lines {
		nodes[i] = Node{
			i,
			line,
		}
	}

	// Build edges weighted by Levenshtein distance between strings
	edges := make([]Edge, 0)
	for i, node1 := range nodes {
		for j, node2 := range nodes {
			if node1.id == node2.id {
				continue
			}
			distance := metrics.NewLevenshtein().Distance(node1.value, node2.value)
			edges = append(edges, Edge{distance, &nodes[i], &nodes[j]})
		}
	}

	// Find minimum spanning tree
	mst := Kruskal(len(nodes), edges)

	groups := extractGroups(mst)

	return groups
}
