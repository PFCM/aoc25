// binary four solves the fourth puzzle.
package main

import (
	"bufio"
	"io"
	"iter"
	"log"
	"os"

	"github.com/pfcm/aoc25"
)

func main() {
	cells, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	aoc25.PrintTiming("Part one", func() int { return partOne(cells) })
}

func partTwo(cells [][]bool) int {
	return 0
}

func partOne(cells [][]bool) int {
	total := 0
	for r, row := range cells {
		for c, cell := range row {
			if cell && accessible(cells, r, c) {
				total++
			}
		}
	}
	return total
}

func accessible(cells [][]bool, i, j int) bool {
	count := 0
	for c := range neighbourhood(cells, i, j) {
		if c {
			count++
		}
	}
	return count < 4
}

func neighbourhood(cells [][]bool, row, col int) iter.Seq[bool] {
	return func(yield func(bool) bool) {
		get := func(i, j int) bool {
			if i < 0 || j < 0 {
				return false
			}
			if i >= len(cells) || j >= len(cells[0]) {
				return false
			}
			return cells[i][j]
		}
		for _, dr := range []int{-1, 0, 1} {
			for _, dc := range []int{-1, 0, 1} {
				if dr == 0 && dc == 0 {
					continue
				}
				if !yield(get(row+dr, col+dc)) {
					return
				}
			}
		}
	}
}

func read(r io.Reader) ([][]bool, error) {
	var (
		scan  = bufio.NewScanner(r)
		cells [][]bool
	)
	for scan.Scan() {
		line := scan.Bytes()
		row := make([]bool, len(line))
		for i, b := range line {
			if b == '@' {
				row[i] = true
			}
		}
		cells = append(cells, row)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return cells, nil
}
