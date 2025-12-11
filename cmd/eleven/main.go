// eleven is the 11th and penultimate day.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pfcm/aoc25"
)

func main() {
	graph, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	aoc25.PrintTiming("Part one", func() int { return partOne(graph) })
}

func partOne(d *devices) int {
	// First do a topological sort starting at "you"
	start := -1
	for i, n := range d.names {
		if n == "you" {
			start = i
			break
		}
	}
	if start == -1 {
		panic("where are you?")
	}

	inDegrees := make([]int, len(d.names))
	for _, es := range d.edges {
		for _, to := range es {
			inDegrees[to]++
		}
	}
	var q []int
	for i, deg := range inDegrees {
		if deg == 0 {
			q = append(q, i)
		}
	}
	sorted := make([]int, 0, len(d.names))
	for len(q) != 0 {
		n := len(q) - 1
		c := q[n]
		q = q[:n]
		sorted = append(sorted, c)

		for _, to := range d.edges[c] {
			inDegrees[to]--
			if inDegrees[to] == 0 {
				q = append(q, to)
			}
		}
	}

	// Now we have them in order, we run through and build a running count
	// of how many ways there are to get to each node.
	paths := make([]int, len(d.names))
	paths[start] = 1
	for _, node := range sorted {
		for _, neighbour := range d.edges[node] {
			paths[neighbour] += paths[node]
		}
	}

	return paths[len(paths)-1]
}

type devices struct {
	names []string
	edges [][]int
}

func read(r io.Reader) (*devices, error) {
	var (
		names       []string
		edgeNames   [][]string
		nameToIndex = make(map[string]int)
		scan        = bufio.NewScanner(r)
	)
	for scan.Scan() {
		l := scan.Text()

		name, edges, ok := strings.Cut(l, ":")
		if !ok {
			return nil, fmt.Errorf("invalid line: %q", l)
		}
		i := len(names)
		names = append(names, name)
		nameToIndex[name] = i
		edgeNames = append(edgeNames, strings.Split(strings.TrimSpace(edges), " "))
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}

	// "out" is special, there are no edges leaving it so we need to add it
	// explicitly.
	nameToIndex["out"] = len(names)
	names = append(names, "out")

	edges := make([][]int, len(names))
	for from, ens := range edgeNames {
		es := make([]int, 0, len(ens))
		for _, en := range ens {
			to, ok := nameToIndex[en]
			if !ok {
				return nil, fmt.Errorf("unknown label in edge: %q", en)
			}
			es = append(es, to)
		}
		edges[from] = es
	}
	return &devices{
		names: names,
		edges: edges,
	}, nil
}
