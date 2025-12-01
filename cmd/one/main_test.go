package main

import (
	"os"
	"testing"
)

func BenchmarkParts(b *testing.B) {
	f, err := os.Open("./inputs/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	input, err := read(f)
	if err != nil {
		b.Fatal(err)
	}
	for _, c := range []struct {
		name string
		f    func([]int) int
	}{{
		name: "one",
		f:    partOne,
	}, {
		name: "two",
		f:    partTwo,
	}} {
		b.Run(c.name, func(b *testing.B) {
			for b.Loop() {
				x := c.f(input)
				_ = x
			}
		})
	}
}
