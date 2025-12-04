// binary four solves the fourth puzzle.
package main

import (
	"bufio"
	"flag"
	"io"
	"iter"
	"log"
	"os"
	"runtime/pprof"

	"github.com/pfcm/aoc25"
)

var (
	profileFlag = flag.String("profile", "", "`path` to write profiles for part two")
)

func main() {
	flag.Parse()

	cells, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	aoc25.PrintTiming("Part one", func() int { return partOne(cells) })

	if *profileFlag != "" {
		finish, err := startProfiles(*profileFlag)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := finish(); err != nil {
				log.Fatal(err)
			}
		}()
	}
	aoc25.PrintTiming("Part two", func() int { return partTwo(cells) })
}

func startProfiles(path string) (func() error, error) {
	// TODO: memory profiles would be cool
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		return nil, err
	}
	return func() error {
		pprof.StopCPUProfile()
		return f.Close()
	}, nil
}

func partTwo(cells [][]bool) int {
	cpy := func(dest, src [][]bool) {
		if len(dest) != len(src) {
			panic("they are different")
		}
		for i := range dest {
			copy(dest[i], src[i])
		}
	}
	next := make([][]bool, len(cells))
	for i := range cells {
		next[i] = make([]bool, len(cells[i]))
	}
	cpy(next, cells)
	var (
		total int
		round = 1
	)
	for round != 0 {
		round = 0

		for r, row := range cells {
			for c, cell := range row {
				if cell && accessible(cells, r, c) {
					round++
					next[r][c] = false
				}
			}
		}
		cpy(cells, next)
		total += round
	}
	return total
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
