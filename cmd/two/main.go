// binary two is day two.
package main

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"slices"
	"strconv"

	"github.com/pfcm/aoc25"
	"github.com/pfcm/it"
)

func main() {
	ranges, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	aoc25.PrintTiming("Part one", func() uint64 { return partOne(ranges) })
	aoc25.PrintTiming("Part two", func() uint64 { return partTwo(ranges) })
}

func partOne(ranges []Range) uint64 {
	valid := func(b []byte) bool {
		if len(b)%2 == 1 {
			// odd length strings can't be made of a repeat
			return true
		}
		n := len(b) / 2
		return !bytes.Equal(b[:n], b[n:])
	}
	sum := uint64(0)
	for num, i := range iterRanges(ranges) {
		if !valid(num) {
			sum += i
		}
	}
	return sum
}

func partTwo(ranges []Range) uint64 {
	valid := func(b []byte) bool {
		for d := 2; len(b)/d > 0; d++ {
			if len(b)%d != 0 {
				continue
			}
			var (
				n   = len(b) / d
				seq = b[:n]
			)
			if it.All(it.Map(it.Batch(slices.Values(b), n), func(s []byte) bool {
				return bytes.Equal(s, seq)
			})) {
				return false
			}
		}
		return true
	}
	sum := uint64(0)
	for num, i := range iterRanges(ranges) {
		if !valid(num) {
			sum += i
		}
	}
	return sum
}

func iterRanges(ranges []Range) iter.Seq2[[]byte, uint64] {
	return func(yield func([]byte, uint64) bool) {
		var scratch []byte
		for _, r := range ranges {
			for i := r.a; i <= r.b; i++ {
				scratch = strconv.AppendUint(scratch[:0], i, 10)
				if !yield(scratch, i) {
					return
				}
			}
		}
	}
}

type Range struct {
	a, b uint64
}

func read(r io.Reader) ([]Range, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	raw = bytes.TrimSpace(raw)
	var ranges []Range
	for pair := range bytes.SplitSeq(raw, []byte{','}) {
		rawA, rawB, ok := bytes.Cut(pair, []byte{'-'})
		if !ok {
			return nil, fmt.Errorf("invalid range: %q", pair)
		}
		a, err := strconv.ParseUint(string(rawA), 10, 64)
		if err != nil {
			return nil, err
		}
		b, err := strconv.ParseUint(string(rawB), 10, 64)
		if err != nil {
			return nil, err
		}
		ranges = append(ranges, Range{a: a, b: b})
	}
	return ranges, nil
}
