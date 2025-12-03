// package aoc25 holds answers to the 2025 advent of code.
package aoc25

import (
	"fmt"
	"time"
)

func PrintTiming[T any](name string, f func() T) {
	t0 := time.Now()
	fmt.Printf("%s: %v (%v)\n", name, f(), time.Since(t0))
}
