package similarity

import (
	"sort"
)

func sortedConnections(nodes []Node) {
	// Sort list of edges by weight
	connections := make([]Connection, 0)
	for _, node := range nodes {
		for _, connection := range node.connections {
			connections = append(connections, connection)
		}
	}

	// Sort list of connections by their weight field
	sort.Slice(connections, func(i, j int) bool {
		return connections[i].weight < connections[j].weight
	})
}

func Kruskal(nodes []Node) []Node {
	connections := sortedConnections(nodes)

	return []Node{}
}
