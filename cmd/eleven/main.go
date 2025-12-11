// eleven is the 11th and penultimate day.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pfcm/aoc25"
)

var (
	oneFlag = flag.Bool("one", true, "whether or not to do part 1")
)

func main() {
	flag.Parse()

	graph, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if *oneFlag {
		aoc25.PrintTiming("Part one", func() int { return partOne(graph) })
	}
	aoc25.PrintTiming("Part two", func() int { return partTwo(graph) })
}

func partTwo(d *devices) int {
	// In the input there's 318388105768112962 paths between
	// svr and out, so we need to be a little bit clever.
	sorted := topo(d)

	find := func(name string) int {
		i := -1
		for j, n := range d.names {
			if n == name {
				i = j
				break
			}
		}
		if i == -1 {
			panic(fmt.Sprintf("no %q", name))
		}
		return i
	}

	svr := find("svr")
	dac := find("dac")
	fft := find("fft")

	paths := make([]int, len(d.names))
	between := func(from, to int) int {
		for i := range paths {
			paths[i] = 0
		}
		paths[from] = 1
		for _, node := range sorted {
			for _, neighbour := range d.edges[node] {
				paths[neighbour] += paths[node]
			}
		}
		return paths[to]
	}
	s2d := between(svr, dac)
	d2f := between(dac, fft)

	s2f := between(svr, fft)
	f2d := between(fft, dac)

	d2o := between(dac, len(d.names)-1)
	f2o := between(fft, len(d.names)-1)

	// fmt.Printf("s2d: %d, d2f: %d, s2f: %d, f2d: %d, d2o: %d, f2o: %d\n", s2d, d2f, s2f, f2d, d2o, f2o)
	return (s2d * d2f * f2o) + (s2f * f2d * d2o)
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
	sorted := topo(d)
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

func topo(d *devices) []int {
	// Kahn's algorithm for topological sort.
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
	return sorted
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
