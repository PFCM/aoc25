// seven is the answers for the seventh of December.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pfcm/aoc25"
)

func main() {
	grid, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part one: %d\n", partOne(grid))
	aoc25.PrintTiming("Part two", func() int { return partTwo(grid) })
}

func partOne(grid [][]Cell) int {
	grid = copyGrid(grid)

	start := -1
	for i, c := range grid[0] {
		if c == Start {
			start = i
			break
		}
	}
	if start == -1 {
		panic("no start position")
	}

	grid[0][start] = Ray

	printGrid(grid)

	splits := 0
	for i, row := range grid[:len(grid)-1] {
		for j, c := range row {
			switch c {
			case Ray:
				switch c2 := grid[i+1][j]; c2 {
				case Start:
					panic("oh my")
				case Ray:
					// this is fine, a splitter probably did
					// the thing through a splitter.
				case Empty:
					grid[i+1][j] = Ray
				case Splitter:
					if j == 0 || j >= len(grid[i+1]) {
						panic("splitter too close to the edge")
					}
					splits++
					// Check what's there before overwriting?
					grid[i+1][j-1] = Ray
					grid[i+1][j+1] = Ray
				default:
					panic(fmt.Errorf("what is %v", c2))
				}
			case Start:
				panic("unpossible")
			case Splitter:
			case Empty:
			default:
			}
		}
		printGrid(grid)
	}
	return splits
}

func partTwo(grid [][]Cell) int {
	// nb. just enumerating all the paths with a graph search was indeed too
	// slow. But just doing the same thing and memoising might be cool.
	memoed := make([][]int, len(grid))
	for i := range memoed {
		memoed[i] = make([]int, len(grid[i]))
		for j := range memoed[i] {
			memoed[i][j] = -1
		}
	}

	var do func(i, j int) int
	do = func(i, j int) int {
		if n := memoed[i][j]; n != -1 {
			return n
		}
		if i == len(grid)-1 {
			memoed[i][j] = 1
			return 1
		}
		// Try and go down.
		switch grid[i+1][j] {
		case Empty:
			memoed[i][j] = do(i+1, j)
			return memoed[i][j]
		case Splitter:
			memoed[i][j] = do(i+1, j-1) + do(i+1, j+1)
			return memoed[i][j]
		}
		panic("sad time")
	}

	start := -1
	for i, c := range grid[0] {
		if c == Start {
			start = i
			break
		}
	}
	if start == -1 {
		panic("no start")
	}

	return do(1, start)
}

func copyGrid(grid [][]Cell) [][]Cell {
	g := make([][]Cell, len(grid))
	for i := range g {
		g[i] = make([]Cell, len(grid[i]))
		copy(g[i], grid[i])
	}
	return g
}

func read(r io.Reader) ([][]Cell, error) {
	var (
		cells [][]Cell
		scan  = bufio.NewScanner(r)
	)
	for scan.Scan() {
		var row []Cell
		for _, r := range scan.Text() {
			c, ok := RuneToCell[r]
			if !ok {
				return nil, fmt.Errorf("unexpected cell %q", r)
			}
			row = append(row, c)
		}
		cells = append(cells, row)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return cells, nil
}

func printGrid(grid [][]Cell) {
	var sb strings.Builder
	for _, row := range grid {
		for _, c := range row {
			sb.WriteString(c.String())
		}
		sb.WriteByte('\n')
	}
	fmt.Println(sb.String())
}

type Cell uint8

const (
	Empty Cell = iota
	Start
	Splitter
	Ray
)

var RuneToCell = map[rune]Cell{
	'.': Empty,
	'S': Start,
	'^': Splitter,
	'|': Ray,
}

var CellToRune = []Cell{
	Empty:    '.',
	Start:    'S',
	Splitter: '^',
	Ray:      '|',
}

func (c Cell) String() string {
	if int(c) >= len(CellToRune) {
		return "?"
	}
	return string(CellToRune[c])
}
