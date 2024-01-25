package similarity

import (
	"fmt"
	"strings"

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
	best := struct {
		node        Node
		totalWeight int
	}{}
	for _, edge := range edges {
		node1Connections := GetNodeConnections(*edge.node1, edges)
		if node1Connections > best.totalWeight {
			best.totalWeight = node1Connections
			best.node = *edge.node1
		}
		node2Connections := GetNodeConnections(*edge.node2, edges)
		if node2Connections > best.totalWeight {
			best.totalWeight = node2Connections
			best.node = *edge.node2
		}
	}

	seen := make(map[int]string)
	extractNodeGroup(best.node, edges, seen, best.totalWeight)

	return []Group{}
}

func extractNodeGroup(node Node, graph Graph, seen map[int]string, weight int) Group {
	edges := GetNodeEdges(node, graph)
	for _, edge := range edges {
		// If node has not been visited yet and the edge will take us to a more similar string
		if _, ok := seen[edge.node1.id]; !ok && edge.weight < weight {
			seen[edge.node1.id] = edge.node1.value
			extractNodeGroup(*edge.node2, graph, seen, edge.weight)
		}
		if _, ok := seen[edge.node2.id]; !ok && edge.weight < weight {
			seen[edge.node2.id] = edge.node2.value
			extractNodeGroup(*edge.node1, graph, seen, edge.weight)
		}
	}

	//convert seen to group
	group := Group{}
	for id := range seen {
		group.Members = append(group.Members, seen[id])
	}
	return group
}

func GetNode(id int, graph Graph) (Node, error) {
	for _, edge := range graph {
		if edge.node1.id == id {
			return *edge.node1, nil
		}
		if edge.node2.id == id {
			return *edge.node2, nil
		}
	}

	return Node{}, fmt.Errorf("Node with id %d not found", id)
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

func GetNodeEdges(node Node, graph Graph) []Edge {
	edges := make([]Edge, 0)
	for _, edge := range graph {
		if edge.node1.id == node.id || edge.node2.id == node.id {
			edges = append(edges, edge)
		}
	}

	return edges
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

func splitLines(text string) []string {
	return strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
}

func FindGroups(logtext string) []Group {
	lines := splitLines(logtext)
	return FindGroupsFromLines(lines)
}

func FindGroupsFromLines(lines []string) []Group {
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
