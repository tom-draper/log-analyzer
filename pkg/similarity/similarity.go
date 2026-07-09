package similarity

import (
	"fmt"
	"sort"
	"strings"

	"github.com/adrg/strutil/metrics"
)

// cutRatio is the normalized edit-distance (distance / longest string) above
// which two lines are considered too dissimilar to belong to the same group.
// Cutting the minimum spanning tree at this threshold performs single-linkage
// clustering: the remaining connected components are the groups.
const cutRatio = 0.6

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

// extractGroups partitions the nodes into groups of similar lines by keeping
// the minimum-spanning-tree edges that join sufficiently similar lines and
// treating the resulting connected components as groups. Nodes that are not
// joined to any other (because every incident edge was cut) form their own
// singleton group. Groups are returned largest first.
func extractGroups(nodes []Node, mst Graph) []Group {
	idToIndex := make(map[int]int, len(nodes))
	for index, node := range nodes {
		idToIndex[node.id] = index
	}

	parent := make([]int, len(nodes))
	for i := range parent {
		parent[i] = i
	}
	var find func(int) int
	find = func(i int) int {
		for parent[i] != i {
			parent[i] = parent[parent[i]] // path halving
			i = parent[i]
		}
		return i
	}
	union := func(a, b int) {
		if ra, rb := find(a), find(b); ra != rb {
			parent[ra] = rb
		}
	}

	for _, edge := range mst {
		if normalizedDistance(edge) <= cutRatio {
			union(idToIndex[edge.node1.id], idToIndex[edge.node2.id])
		}
	}

	membersByRoot := make(map[int][]string)
	for index, node := range nodes {
		root := find(index)
		membersByRoot[root] = append(membersByRoot[root], node.value)
	}

	groups := make([]Group, 0, len(membersByRoot))
	for _, members := range membersByRoot {
		groups = append(groups, Group{Members: members})
	}
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i].Members) > len(groups[j].Members)
	})
	return groups
}

// normalizedDistance returns an edge's edit distance scaled by the length of
// the longer of its two strings, giving a 0..1 dissimilarity that is comparable
// across lines of different lengths.
func normalizedDistance(edge Edge) float64 {
	longest := len(edge.node1.value)
	if len(edge.node2.value) > longest {
		longest = len(edge.node2.value)
	}
	if longest == 0 {
		return 0
	}
	return float64(edge.weight) / float64(longest)
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

	// Build edges weighted by Levenshtein distance between strings. The
	// distance is symmetric and all consumers treat edges as undirected, so
	// only the unique pairs (j > i) are needed.
	lev := metrics.NewLevenshtein()
	edges := make([]Edge, 0, len(nodes)*(len(nodes)-1)/2)
	for i := range nodes {
		for j := i + 1; j < len(nodes); j++ {
			distance := lev.Distance(nodes[i].value, nodes[j].value)
			edges = append(edges, Edge{distance, &nodes[i], &nodes[j]})
		}
	}

	// Find minimum spanning tree
	mst := Kruskal(len(nodes), edges)

	groups := extractGroups(nodes, mst)

	return groups
}
