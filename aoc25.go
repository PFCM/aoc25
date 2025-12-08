// package aoc25 holds answers to the 2025 advent of code.
package aoc25

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

func PrintTiming[T any](name string, f func() T) {
	t0 := time.Now()
	fmt.Printf("%s: %v (%v)\n", name, f(), time.Since(t0))
}

type IntVector[S constraints.Signed] struct {
	X, Y, Z S
}

func NewIntVectorFromString[S constraints.Signed](s string) (IntVector[S], error) {
	pieces := strings.Split(s, ",")
	if l := len(pieces); l != 3 {
		return IntVector[S]{}, fmt.Errorf("invalid vector: need 3 components: got %q", s)
	}
	var numbers [3]S
	for i, l := range pieces {
		n, err := strconv.ParseInt(l, 10, 64)
		if err != nil {
			return IntVector[S]{}, nil
		}
		numbers[i] = S(n)
	}
	return IntVector[S]{
		X: numbers[0],
		Y: numbers[1],
		Z: numbers[2],
	}, nil
}

func (iv IntVector[S]) EuclideanDistance(other IntVector[S]) float64 {
	if iv == other {
		return 0
	}
	d := IntVector[S]{iv.X - other.X, iv.Y - other.Y, iv.Z - other.Z}
	return math.Sqrt(float64(
		d.X*d.X + d.Y*d.Y + d.Z*d.Z,
	))
}
