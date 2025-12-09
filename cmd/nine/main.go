// nine is day nine.
package main

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

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

func partTwo(points []point) int64 {
	inShape := func(p point) bool {
		// tests if the point is in the shape, using a crossing number
		// algorithm like those described in
		// https://web.archive.org/web/20130126163405/http://geomalgorithms.com/a03-_inclusion.html
		// We can be a bit cheeky though because it's all right angles.
		// The idea here is we're projecting a horizontal ray to the right
		// of p, and counting the number of times it crosses an edge.
		cn := 0
		for i := range points {
			start, end := points[i], points[(i+1)%len(points)]
			if p.y == start.y && p.y == end.y {
				// both straight horizontal lines
				if p.x < start.x {
					cn++
				}
				if p.x < end.x {
					cn++
				}
				continue
			}
			if (p.y >= start.y && p.y <= end.y) || (p.y >= end.y && p.y <= start.y) {
				// By construction, start.x and end.x must be the
				// same, and the intersection with the ray is at
				// (p.y, start.x). So all we need to know is if
				// p is to the left of the line.
				if p.x < start.x {
					cn++
				}
			}
		}
		return (cn % 2) == 1
	}
	var rp [4]point
	contained := func(p, q point) bool {
		// Returns true iff the rectangle defined by p and q is entirely
		// contained with the shape described by lines.
		minX, minY := min(p.x, q.x), min(p.y, q.y)
		maxX, maxY := max(p.x, q.x), max(p.y, q.y)
		rp[0] = point{x: minX, y: minY} // top left
		rp[1] = point{x: maxX, y: minY} // top right
		rp[2] = point{x: maxX, y: maxY} // bottom right
		rp[3] = point{x: minX, y: maxY} // bottom left
		// To avoid getting a situation like
		// ########
		// .......#
		// ....####
		// ....#...
		// ....#...
		// ....####
		// .......#
		// .......#
		// erroneously allowing a rectangle that covers the cutout,
		// check every point along the perimeter.
		// Which is a lot of points.
		// There's probably a better way involving breaking each line
		// segment into chunks according to where it intersects
		// the shape, but even that doesn't seem like it'd be great.
		for i, start := range rp {
			end := rp[(i+1)%4]
			for p := range iterLine(start, end) {
				if !inShape(p) {
					return false
				}
			}
		}
		return true
	}

	checked := int64(0)
	total := len(points) * len(points) / 2
	workers := runtime.GOMAXPROCS(0)
	results := make([][]int64, workers)
	ixChan := make(chan int)
	var g sync.WaitGroup
	for i := range workers {
		go func() {
			defer g.Done()
			for j := range ixChan {
				p := points[j]
				for _, q := range points[j+1:] {
					if contained(p, q) {
						results[i] = append(results[i], area(p, q))
					}
					if c := atomic.AddInt64(&checked, 1); c%1000 == 0 {

						fmt.Printf("\r%d/%d", checked, total)
					}
				}
			}
		}()
		g.Add(1)
	}
	for i := range points {
		ixChan <- i
	}
	close(ixChan)
	g.Wait()
	fmt.Println()

	largest := int64(0)
	for _, rs := range results {
		for _, r := range rs {
			if r > largest {
				largest = r
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
