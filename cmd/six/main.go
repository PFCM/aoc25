// six is the answer to day six.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	p1, err := partOne(bytes.NewReader(input))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part one: %d\n", p1)

	p2, err := partTwo(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part two: %d\n", p2)
}

func partTwo(input []byte) (int, error) {
	ps, err := read2(input)
	if err != nil {
		return 0, err
	}
	return calculateAndSum(ps), nil
}

func partOne(r io.Reader) (int, error) {
	ps, err := read1(r)
	if err != nil {
		return 0, err
	}
	return calculateAndSum(ps), nil
}

func calculateAndSum(ps []problem) int {
	sum := 0
	for _, p := range ps {
		sum += p.Calculate()
	}
	return sum
}

type problem struct {
	inputs []int
	op     op
}

type op uint8

const (
	opAdd op = '+'
	opMul op = '*'
)

func (p problem) Calculate() int {
	var f func(a, b int) int
	switch p.op {
	case opAdd:
		f = func(a, b int) int { return a + b }
	case opMul:
		f = func(a, b int) int { return a * b }
	default:
		panic(":/")
	}
	x := p.inputs[0]
	for _, y := range p.inputs[1:] {
		x = f(x, y)
	}
	return x
}

func read2(raw []byte) ([]problem, error) {
	lines := bytes.Split(raw, []byte{'\n'})
	if n := len(lines) - 1; len(lines[n]) == 0 {
		lines = lines[:n]
	}

	lastLine := lines[len(lines)-1]
	lines = lines[:len(lines)-1]

	newLines := make([][]byte, len(lines[0]))
	for i := range newLines {
		newLines[i] = make([]byte, len(lines))
	}
	for i := range lines {
		for j := range lines[i] {
			newLines[j][i] = lines[i][j]
		}
	}

	var (
		p  problem
		ps []problem
	)
	for _, l := range newLines {
		l = bytes.TrimSpace(l)
		if len(l) == 0 {
			ps = append(ps, p)
			p = problem{}
			continue
		}
		n, err := strconv.Atoi(string(l))
		if err != nil {
			return nil, err
		}
		p.inputs = append(p.inputs, n)
	}
	ps = append(ps, p)

	i := 0
	for _, b := range lastLine {
		if b == ' ' {
			continue
		}
		switch b {
		case '+':
			ps[i].op = opAdd
		case '*':
			ps[i].op = opMul
		default:
			return nil, fmt.Errorf("unknown op %q", b)
		}
		i++
	}
	return ps, nil
}

func read1(r io.Reader) ([]problem, error) {
	var (
		problems []problem
		scan     = bufio.NewScanner(r)
	)
	for first := true; scan.Scan(); first = false {
		pieces := strings.Split(scan.Text(), " ")
		w := 0
		for _, p := range pieces {
			if p == "" {
				continue
			}
			pieces[w] = p
			w++
		}
		pieces = pieces[:w]

		if p := pieces[0]; p == "*" || p == "+" {
			for i, p := range pieces {
				var o op
				switch p {
				case "*":
					o = opMul
				case "+":
					o = opAdd
				default:
					return nil, fmt.Errorf("unknown opeartion %q", p)
				}
				problems[i].op = o
			}
			continue // maybe break
		}

		for i, p := range pieces {
			n, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			if first {
				problems = append(problems, problem{
					inputs: []int{n},
				})
			} else {
				problems[i].inputs = append(problems[i].inputs, n)
			}
		}

	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return problems, nil
}
