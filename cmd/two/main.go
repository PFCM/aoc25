// binary two is day two.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	ranges, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part one: %d\n", partOne(ranges))
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
	var (
		sum     uint64 = 0
		scratch []byte
	)
	for _, r := range ranges {
		for i := r.a; i <= r.b; i++ {
			scratch = strconv.AppendUint(scratch[:0], i, 10)
			if !valid(scratch) {
				sum += i
			}
		}
	}
	return sum
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
