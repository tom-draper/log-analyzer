package similarity

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/adrg/strutil/metrics"
)

type Line struct {
	lineIndex int
	line      string
}

type Group struct {
	examples []string
}

func getStartingCentroids(lines []Line, k int) []Line {
	linesCopy := append(lines[:0:0], lines...)
	samples := make([]Line, k)

	for i := 0; i < k; i++ {
		r := int(rand.Int63n(int64(len(linesCopy))))
		samples[i] = linesCopy[r]

		// remove the sample from the copy slice
		linesCopy[r], linesCopy[len(linesCopy)-1] = linesCopy[len(linesCopy)-1], linesCopy[r]
		linesCopy = linesCopy[:len(linesCopy)-1]
	}

	return samples
}

func getNearestCluster(line Line, clusters []Line) struct {
	distance float64
	index    int
	line     Line
} {
	best := struct {
		distance float64
		index    int
		line     Line
	}{distance: math.MaxFloat64}
	for i, cluster := range clusters {
		var distance float64 = metrics.NewLevenshtein().Distance(line.line, cluster.line)
		fmt.Println(distance)
		if distance < best.distance {
			best = struct {
				distance float64
				index    int
				line     Line
			}{distance, i, line}
		}
	}
	return best
}

func buildLines(strings []string) []Line {
	lines := make([]Line, len(strings))
	for i, line := range strings {
		lines[i] = Line{
			i,
			line,
		}
	}
	return lines
}

type Node struct {
	id          int
	value       string
	connections []Connection
}

type Connection struct {
	weight int
	node   *Node
}

func FindGroups(lines []string) []Group {
	// lines := buildLines(strings)

	// k := 3
	// centroids := getStartingCentroids(lines, k)

	// groups := make([][]Line, k)
	// for _, line := range lines {
	// 	nearest := getNearestCluster(line, centroids)
	// 	groups[nearest.index] = append(groups[nearest.index], nearest.line)
	// }
	// return []Group{}

	nodes := make([]Node, len(lines))
	for i, line := range lines {
		// create new node
		node := Node{
			i,
			line,
			[]Connection{},
		}
		nodes[i] = node
	}

	for i, line1 := range lines {
		for j, line2 := range lines {
			distance := metrics.NewLevenshtein().Distance(line1, line2)
			nodes[i].connections = append(nodes[i].connections, Connection{distance, &nodes[j]})
		}
	}

	// Find minimum spanning tree
	mst := Kruskal(nodes)

	return []Group
}
