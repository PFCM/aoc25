// binary one solves the puzzle for the first of December.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	input, err := read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("part one: %d\n", partOne(input))
}

func partOne(turns []int) int {
	var (
		pos   = 50
		zeros = 0
	)
	for _, t := range turns {
		pos = (pos + t) % 100
		if pos == 0 {
			zeros++
		}
	}
	return zeros
}

// read reads the input from stdin, as a list of integers: right rotations are
// positive, left negative.
func read() ([]int, error) {
	var (
		scan    = bufio.NewScanner(os.Stdin)
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
