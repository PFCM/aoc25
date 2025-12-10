// ten is the tenth day
package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/pfcm/it"

	"github.com/pfcm/aoc25"
)

var workersFlag = flag.Int("workers", 10, "parallelism for part 2")

func main() {
	flag.Parse()

	machines, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	aoc25.PrintTiming("Part one", func() int { return partOne(machines) })
	aoc25.PrintTiming("Part two", func() int32 { return partTwo(machines) })
}

func partOne(ms []machine) int {
	n := 0
	for _, m := range ms {
		n += m.turnOn()
	}
	return n
}

func partTwo(ms []machine) int32 {
	var (
		n     atomic.Int32
		done  atomic.Int32
		g     sync.WaitGroup
		batch = max(1, len(ms)/(*workersFlag))
	)
	for b := range it.Batch(slices.Values(ms), batch) {
		b := slices.Clone(b)
		g.Go(func() {
			for _, m := range b {
				v := m.setJoltages()
				n.Add(int32(v))
				fmt.Printf("%d/%d done\n", done.Add(1), len(ms))
			}
		})
	}
	g.Wait()
	return n.Load()
}

func read(r io.Reader) ([]machine, error) {
	var (
		machines []machine
		scan     = bufio.NewScanner(r)
	)
	for scan.Scan() {
		m, err := fromBytes(scan.Bytes())
		if err != nil {
			return nil, err
		}
		machines = append(machines, m)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return machines, nil
}

type machine struct {
	targetLights uint16   // 0 bit means we want it off, 1 on
	buttons      []uint16 // 0 bit means no change, 1 means it toggles
	joltages     []int    // actually numbers
}

func fromBytes(b []byte) (machine, error) {
	pieces := bytes.Split(b, []byte{' '})
	if len(pieces) < 3 {
		return machine{}, fmt.Errorf("invalid machine: %q", b)
	}
	lights := uint16(0)
	for i, l := range pieces[0][1 : len(pieces[0])-1] {
		if l == '#' {
			lights |= 1 << i
		}
	}
	var buttons []uint16
	for _, button := range pieces[1 : len(pieces)-1] {
		butt := uint16(0)
		ns, err := numbers(button)
		if err != nil {
			return machine{}, err
		}
		for _, n := range ns {
			butt |= 1 << n
		}
		buttons = append(buttons, butt)
	}

	js, err := numbers(pieces[len(pieces)-1])
	if err != nil {
		return machine{}, err
	}
	m := machine{
		targetLights: lights,
		buttons:      buttons,
		joltages:     js,
	}
	return m, nil
}

func (m machine) String() string {
	return fmt.Sprintf("[%08b] %08b %v", m.targetLights, m.buttons, m.joltages)
}

func (m machine) turnOn() int {
	// returns the shortest number of button presses to turn the machine on:
	// to make an entirely off set of lights match m.targetLights.
	type step struct {
		length int
		lights uint16
	}
	todo := []step{{length: 0, lights: 0}}
	for range 100000000 {
		s := todo[0]
		todo = todo[1:]
		if s.lights == m.targetLights {
			return s.length
		}
		for _, b := range m.buttons {
			todo = append(todo, step{
				length: s.length + 1,
				lights: s.lights ^ b,
			})
		}
	}
	panic("max iterations")
}

func (m machine) setJoltages() int {
	match := func(js []int) bool {
		return it.All(it.Map2x1(
			it.Zip(slices.Values(js), slices.Values(m.joltages)),
			func(a, b int) bool {
				return a == b
			}))
	}
	press := func(js []int, button uint16) []int {
		// RIP
		js = slices.Clone(js)
		for i := range js {
			if (button>>i)&1 == 1 {
				js[i]++
			}
		}
		return js
	}
	heuristic := func(js []int) float64 {
		// TODO: does squared euclidean make any sense at all? can I even remember
		// the requirements for a valid heuristic?
		// var d2 float64
		// for a, b := range it.Zip(slices.Values(js), slices.Values(m.joltages)) {
		// 	d := float64(a - b)
		// 	d2 += d * d
		// }
		// return d2
		// TODO: does l1 even make any sense at all? can I even remember
		// the requirements for a valid heuristic?
		var d2 int
		for a, b := range it.Zip(slices.Values(js), slices.Values(m.joltages)) {
			d := (a - b)
			if d < 0 {
				d = -d
			}
			d2 += d
		}
		return float64(d2)
	}
	// lol
	visited := make(map[string]int)
	todo := jheap{{length: 0, joltages: make([]int, len(m.joltages))}}
	for range 1000000000 {
		s := heap.Pop(&todo).(jnode)

		key := fmt.Sprint(s.joltages)
		if n, ok := visited[key]; ok && n <= s.length {
			// We've previous found a way to this permutation that
			// was better (or the same, either way).
			continue
		}
		visited[key] = s.length
		if match(s.joltages) {
			return s.length
		}
		for _, b := range m.buttons {
			js := press(s.joltages, b)
			skip := false
			for x, y := range it.Zip(slices.Values(js), slices.Values(m.joltages)) {
				if x > y {
					// we can't go down, there's no point
					// pushing this.
					skip = true
					break
				}
			}
			if skip {
				continue
			}
			heap.Push(&todo, jnode{
				h:        heuristic(js),
				length:   s.length + 1,
				joltages: js,
			})
		}
	}
	panic("max iterations")
}

func numbers(b []byte) ([]int, error) {
	var nums []int
	for num := range bytes.SplitSeq(b[1:len(b)-1], []byte{','}) {
		n, err := strconv.Atoi(string(num))
		if err != nil {
			return nil, err
		}
		nums = append(nums, n)
	}
	return nums, nil
}

type jnode struct {
	h        float64
	length   int
	joltages []int
}

type jheap []jnode

func (jh jheap) Len() int           { return len(jh) }
func (jh jheap) Less(i, j int) bool { return jh[i].h < jh[j].h }
func (jh jheap) Swap(i, j int)      { jh[i], jh[j] = jh[j], jh[i] }
func (j *jheap) Push(a any)         { *j = append(*j, a.(jnode)) }

func (j *jheap) Pop() any {
	n := len(*j) - 1
	jn := (*j)[n]
	*j = (*j)[:n]
	return jn
}
