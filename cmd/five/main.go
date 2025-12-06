// five is the solution to day five.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/pfcm/aoc25"
)

func main() {
	ranges, ids, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	aoc25.PrintTiming("Sort input", func() string {
		slices.SortFunc(ranges, func(a, b Range) int {
			return a.compare(b)
		})
		return ":)"
	})

	aoc25.PrintTiming("Part one", func() int { return partOne(ranges, ids) })
	aoc25.PrintTiming("Part two", func() uint64 { return partTwo(ranges) })
}

func partTwo(ranges []Range) uint64 {
	var (
		newRanges []Range
		r         = ranges[0]
	)
	for _, r2 := range ranges[1:] {
		if r.end >= r2.start {
			// The ranges overlap, merge them.
			r.end = max(r.end, r2.end) // max to handle r2 entirely within r
			continue
		}
		newRanges = append(newRanges, r)
		r = r2
	}
	newRanges = append(newRanges, r)

	count := uint64(0)
	for _, r := range newRanges {
		n := r.end - r.start + 1
		count += n
	}
	return count
}

func partOne(ranges []Range, ids []uint64) int {
	count := 0
	for _, id := range ids {
		for _, r := range ranges {
			if r.contains(id) {
				count++
				break
			}
		}
	}
	return count
}

type Range struct {
	start, end uint64 // end is inclusive
}

func (r Range) contains(i uint64) bool {
	return i >= r.start && i <= r.end
}

func (r Range) compare(s Range) int {
	switch {
	case r.start < s.start:
		return -1
	case r.start > s.start:
		return 1
	case r.end < s.end:
		return -1
	case r.end > s.end:
		return 1
	}
	return 0
}

func read(r io.Reader) ([]Range, []uint64, error) {
	var (
		scan   = bufio.NewScanner(r)
		ranges []Range
		ids    []uint64
	)
	for scan.Scan() {
		l := scan.Text()
		if l == "" {
			break
		}
		nums := strings.Split(l, "-")
		if len(nums) != 2 {
			return nil, nil, fmt.Errorf("unexpected input range: %q", l)
		}
		a, err := strconv.ParseUint(nums[0], 10, 64)
		if err != nil {
			return nil, nil, err
		}
		b, err := strconv.ParseUint(nums[1], 10, 64)
		if err != nil {
			return nil, nil, err
		}
		if a > b {
			return nil, nil, fmt.Errorf("invalid range %d-%d", a, b)
		}
		ranges = append(ranges, Range{start: a, end: b})
	}
	if err := scan.Err(); err != nil {
		return nil, nil, err
	}
	for scan.Scan() {
		n, err := strconv.ParseUint(scan.Text(), 10, 64)
		if err != nil {
			return nil, nil, err
		}
		ids = append(ids, n)
	}
	if err := scan.Err(); err != nil {
		return nil, nil, err
	}
	return ranges, ids, nil
}
