// binary one solves the puzzle for the first of December.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/pfcm/aoc25"
)

func main() {
	input, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	aoc25.PrintTiming("Part one", func() int { return partOne(input) })
	aoc25.PrintTiming("Part two", func() int { return partTwo(input) })
}

func partTwo(turns []int) int {
	var (
		pos   = 50
		zeros = 0
	)
	for _, t := range turns {
		// lol
		s := 1
		if t < 0 {
			s = -1
		}
		if t < 0 {
			t = -t
		}
		for range t {
			pos += s
			if pos == 0 {
				zeros++
			}
			switch pos {
			case -1:
				pos = 99
			case 100:
				pos = 0
				zeros++
			}
		}
	}
	return zeros
}

func partOne(turns []int) int {
	clampMod := func(i int) int {
		for i < 0 {
			i += 100
		}
		return i % 100
	}
	var (
		pos   = 50
		zeros = 0
	)
	for _, t := range turns {
		pos = clampMod(pos + t)
		if pos == 0 {
			zeros++
		}
	}
	return zeros
}

// read reads the input from the provided reader, as a list of integers: right
// rotations are positive, left negative.
func read(r io.Reader) ([]int, error) {
	var (
		scan    = bufio.NewScanner(r)
		results []int
	)
	for scan.Scan() {
		line := scan.Text()
		num, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}
		switch line[0] {
		case 'L':
			num = -num
		case 'R':
			// all good
		default:
			return nil, fmt.Errorf("bad line: %q", line)
		}
		results = append(results, num)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
