// eight is a solution for day 8.
package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"maps"
	"os"
	"slices"

	"github.com/pfcm/aoc25"
)

var (
	joins = flag.Int("joins", 1000, "number of lights to join for part 1")
)

func main() {
	flag.Parse()

	inputs, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	aoc25.PrintTiming("Part one", func() int { return partOne(inputs, *joins) })
	aoc25.PrintTiming("Part two", func() int { return partTwo(inputs) })
}

func partOne(vs []aoc25.IntVector[int], joins int) int {
	// 1,000,000 elements for the real input, probably not big
	// enough to be an issue.
	distances := make([]float64, len(vs)*len(vs))
	index := func(i, j int) int {
		return i*len(vs) + j
	}
	fromIndex := func(i int) (int, int) {
		return i % len(vs), i / len(vs)
	}
	// TODO: we only actually need half of these :/
	for i, v1 := range vs {
		for j, v2 := range vs {
			distances[index(i, j)] = v1.EuclideanDistance(v2)
		}
	}
	indexes := make([]int, len(distances))
	for i := range indexes {
		indexes[i] = i
	}

	slices.SortFunc(indexes, func(i, j int) int {
		d1, d2 := distances[i], distances[j]
		if d1 < d2 {
			return -1
		}
		if d1 > d2 {
			return 1
		}
		return 0
	})
	// This is all pretty dubious.
	cut := 0
	for n, i := range indexes {
		d := distances[i]
		x, y := fromIndex(i)
		if d != 0 || x != y {
			cut = n
			break
		}
	}
	indexes = indexes[cut:]
	indexes = slices.CompactFunc(indexes, func(i, j int) bool {
		x1, y1 := fromIndex(i)
		x2, y2 := fromIndex(j)
		return x1 == y2 && y1 == x2
	})

	sets := make(map[aoc25.IntVector[int]]*ufNode, len(vs))
	for _, v := range vs {
		n := &ufNode{value: v, size: 1}
		n.parent = n
		sets[v] = n
	}

	for _, i := range indexes[:joins] {
		x, y := fromIndex(i)
		v1, v2 := vs[x], vs[y]
		sets[v1].Union(sets[v2])
	}

	for v, n := range sets {
		if n.Find() != n {
			delete(sets, v)
		}
	}

	nodes := slices.Collect(maps.Values(sets))
	slices.SortFunc(nodes, func(a, b *ufNode) int {
		return b.size - a.size
	})
	product := 1
	for _, n := range nodes[:3] {
		product *= n.size
	}

	return product
}

func partTwo(vs []aoc25.IntVector[int]) int {
	// TODO: factor out the common bits :(
	distances := make([]float64, len(vs)*len(vs))
	index := func(i, j int) int {
		return i*len(vs) + j
	}
	fromIndex := func(i int) (int, int) {
		return i % len(vs), i / len(vs)
	}
	// TODO: we only actually need half of these :/
	for i, v1 := range vs {
		for j, v2 := range vs {
			distances[index(i, j)] = v1.EuclideanDistance(v2)
		}
	}
	indexes := make([]int, len(distances))
	for i := range indexes {
		indexes[i] = i
	}

	slices.SortFunc(indexes, func(i, j int) int {
		d1, d2 := distances[i], distances[j]
		if d1 < d2 {
			return -1
		}
		if d1 > d2 {
			return 1
		}
		return 0
	})
	// This is all pretty dubious.
	cut := 0
	for n, i := range indexes {
		d := distances[i]
		x, y := fromIndex(i)
		if d != 0 || x != y {
			cut = n
			break
		}
	}
	indexes = indexes[cut:]
	indexes = slices.CompactFunc(indexes, func(i, j int) bool {
		x1, y1 := fromIndex(i)
		x2, y2 := fromIndex(j)
		return x1 == y2 && y1 == x2
	})

	sets := make(map[aoc25.IntVector[int]]*ufNode, len(vs))
	for _, v := range vs {
		n := &ufNode{value: v, size: 1}
		n.parent = n
		sets[v] = n
	}
	// This is where we diverge from part 1: just keep joining the
	// closest ones together until they're all in the same set.
	for _, i := range indexes {
		x, y := fromIndex(i)
		v1, v2 := vs[x], vs[y]
		u := sets[v1].Union(sets[v2])
		if u.size == len(vs) {
			return v1.X * v2.X
		}
	}
	panic("oh no")
}

type ufNode struct {
	parent *ufNode
	value  aoc25.IntVector[int]
	size   int
}

func (u *ufNode) Find() *ufNode {
	for u.parent != u {
		u, u.parent = u.parent, u.parent.parent
	}
	return u
}

func (u *ufNode) Union(v *ufNode) *ufNode {
	u = u.Find()
	v = v.Find()
	if u == v {
		return u
	}
	if u.size < v.size {
		u, v = v, u
	}
	v.parent = u
	u.size += v.size
	return u
}

func read(r io.Reader) ([]aoc25.IntVector[int], error) {
	var (
		results []aoc25.IntVector[int]
		scan    = bufio.NewScanner(r)
	)
	for scan.Scan() {
		v, err := aoc25.NewIntVectorFromString[int](scan.Text())
		if err != nil {
			return nil, err
		}
		results = append(results, v)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
