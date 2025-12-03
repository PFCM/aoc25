// binary three is the answer for the third day.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

func main() {
	banks, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part one: %d\n", partOne(banks))
	fmt.Printf("Part two: %d\n", partTwo(banks))
}

func partTwo(banks [][]uint8) uint64 {
	var (
		sum = uint64(0)
		tmp [12]uint8
	)
	for _, bank := range banks {
		// Same as below, we just have to do it twelve times instead of
		// two.
		start := 0
		for i := range 12 {
			var (
				end      = len(bank) - 11 + i // one after last available digit
				max      uint8
				maxIndex = start
			)
			for j := start; j < end; j++ {
				b := bank[j]
				if b > max {
					max = b
					maxIndex = j
				}
			}
			tmp[i] = max
			start = maxIndex + 1
		}
		inc := uint64(0)
		for i, b := range tmp {
			r := uint64(math.Pow10(11 - i))
			inc += uint64(b) * r
		}
		sum += inc
	}
	return sum
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
