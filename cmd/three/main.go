// binary three is the answer for the third day.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	banks, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part one: %d\n", partOne(banks))
}

func partOne(banks [][]uint8) int {
	sum := 0
	for _, bank := range banks {
		// We're always going to want to start with the highest number
		// (that isn't the last one). The only trick will be what to do
		// when there are duplicates, although in that case we should
		// just take the first.
		var (
			max      uint8
			maxIndex int
		)
		for i, b := range bank[:len(bank)-1] {
			if b > max {
				max = b
				maxIndex = i
			}
		}
		var next uint8
		for _, b := range bank[maxIndex+1:] {
			if b > next {
				next = b
			}
		}
		sum += int(max)*10 + int(next)
	}
	return sum
}

func read(r io.Reader) ([][]uint8, error) {
	var (
		banks [][]uint8
		scan  = bufio.NewScanner(r)
	)
	for scan.Scan() {
		line := scan.Bytes()
		bank := make([]uint8, len(line))
		for i := range bank {
			bank[i] = line[i] - '0'
		}
		banks = append(banks, bank)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return banks, nil
}
