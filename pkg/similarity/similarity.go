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

type Graph []Edge

func (g Graph) String() string {
	s := ""
	for _, edge := range g {
		s += edge.String() + "\n"
	}
	return s
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

func GetNodes(graph Graph) []Node {
	// Set of nodes
	nodes := make(map[int]string)
	for _, edge := range graph {
		nodes[edge.node1.id] = edge.node1.value
		nodes[edge.node2.id] = edge.node2.value
	}

	//convert to slice
	nodesSlice := make([]Node, len(nodes))
	i := 0
	for id, value := range nodes {
		nodesSlice[i] = Node{id, value}
		i++
	}

	return nodesSlice
}

func GetNodeConnections(node Node, graph Graph) int {
	totalWeight := 0
	for _, edge := range graph {
		if edge.node1.id == node.id || edge.node2.id == node.id {
			totalWeight += edge.weight
		}
	}
	return totalWeight
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

	fmt.Println(mst.String())

	groups := extractGroups(mst)

	return groups
}
