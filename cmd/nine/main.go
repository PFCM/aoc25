// nine is day nine.
package main

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pfcm/aoc25"
)

func main() {
	points, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	aoc25.PrintTiming("Part one", func() int64 { return partOne(points) })
	aoc25.PrintTiming("Part two", func() int64 { return partTwo(points) })
}

func partOne(points []point) int64 {
	largest := int64(0)
	for i, p := range points {
		for _, q := range points[i+1:] {
			if a := area(p, q); a > largest {
				largest = a
			}
		}
	}
	return largest
}

/*
 0123456789
0..........
1.###..####
2.#.#..#..#
3.#.####..#
4.#########
*/

func partTwo(points []point) int64 {
	contained := func(p, q point) (result bool) {
		minX, minY := min(p.x, q.x), min(p.y, q.y)
		maxX, maxY := max(p.x, q.x), max(p.y, q.y)
		// If there are any points in the shape that are inside the
		// rectangle that are not on the edge, then the rectangle _must_
		// go outside the shape.
		for _, p := range points {
			if p.x <= minX || p.x >= maxX {
				continue
			}
			if p.y <= minY || p.y >= maxY {
				continue
			}
			return false
		}
		// This is necessary but not sufficient: there could also be
		// lines in the shape that entirely cross the rectangle, which
		// would also mean no dice.
		for i := range points {
			start, end := points[i], points[(i+1)%len(points)]
			// TODO: these conditions seem unreasonably complicated
			if start.x == end.x {
				// vertical line
				if start.x <= minX || start.x >= maxX {
					continue
				}
				start.y, end.y = min(start.y, end.y), max(start.y, end.y)
				if start.y <= minY && end.y >= maxY {
					// crosses the rectangle
					return false
				}
				continue
			} else if start.y == end.y {
				// horizontal line
				if start.y <= minY || start.y >= maxY {
					continue
				}
				start.x, end.x = min(start.x, end.x), max(start.x, end.x)
				if start.x <= minX && end.x >= maxX {
					return false
				}
				continue
			} else {
				panic("impossible?")
			}
		}
		// Maybe this is enough?
		return true
	}

	largest := int64(0)
	for i, p := range points {
		for _, q := range points[i+1:] {
			if a := area(p, q); a > largest && contained(p, q) {
				largest = a
			}

		}
	}
	return largest
}

func iterLine(start, end point) iter.Seq[point] {
	sgn := func(i int64) int64 {
		if i < 0 {
			return -1
		}
		if i > 0 {
			return 1
		}
		return 0
	}
	d := point{sgn(end.x - start.x), sgn(end.y - start.y)}

	return func(yield func(point) bool) {
		for x := start; x != end; x = x.add(d) {
			if !yield(x) {
				return
			}
		}
		yield(end)
	}
}

func area(a, b point) int64 {
	minPoint := point{
		x: min(a.x, b.x),
		y: min(a.y, b.y),
	}
	maxPoint := point{
		x: max(a.x, b.x),
		y: max(a.y, b.y),
	}
	return (maxPoint.x - minPoint.x + 1) * (maxPoint.y - minPoint.y + 1)
}

type point struct {
	x, y int64
}

func (p point) add(q point) point { return point{p.x + q.x, p.y + q.y} }

func read(r io.Reader) ([]point, error) {
	var (
		results []point
		scan    = bufio.NewScanner(r)
	)
	for scan.Scan() {
		nums := strings.Split(scan.Text(), ",")
		if l := len(nums); l != 2 {
			return nil, fmt.Errorf("invalid line: %q", l)
		}
		x, err := strconv.ParseInt(nums[0], 10, 64)
		if err != nil {
			return nil, err
		}
		y, err := strconv.ParseInt(nums[1], 10, 64)
		if err != nil {
			return nil, err
		}
		results = append(results, point{x: x, y: y})
	}
	return results, nil
}
