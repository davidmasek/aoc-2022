package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Screen struct {
	line string
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (s *Screen) draw(cycle, sprite_position int) {
	draw_position := (cycle - 1) % 40
	if abs(draw_position-sprite_position) <= 1 {
		s.line = s.line + "#"
	} else {
		s.line = s.line + "."
	}
	if len(s.line) == 40 {
		fmt.Println(s.line, cycle)
		s.line = ""
	}
}

func main() {
	fmt.Println(len(os.Args), os.Args)
	fmt.Println("-----")
	if len(os.Args) != 2 {
		fmt.Println("Usage: XY.go XY.in")
		os.Exit(1)
	}
	var filename string
	test := true
	if test {
		filename = strings.Replace(os.Args[1], ".ex", ".in", 1)
	} else {
		filename = os.Args[1]
	}
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	i := 0
	// change to be applied after the cycle
	var changes []int = make([]int, 300)
	// value during the cycle
	var values []int = make([]int, 300)
	changes[0] = 1

	screen := Screen{}

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		// noop takes one cycle
		i++
		if line[:4] == "noop" {
			continue
		}
		// addx takes two cycles
		var v int
		fmt.Sscanf(line, "addx %d", &v)
		i++
		changes[i] = v
	}

	sum := 0
	for i := 1; i < len(changes); i++ {
		sum += changes[i-1]
		values[i] = sum
	}
	// for j := 0; j < 10; j++ {
	// 	fmt.Println(j, changes[j], values[j])
	// }
	fmt.Println(values[20], values[60], values[100], values[140], values[180], values[220])
	fmt.Println(
		values[20]*20 +
			values[60]*60 +
			values[100]*100 +
			values[140]*140 +
			values[180]*180 +
			values[220]*220)
	fmt.Println("-----")
	for i := 1; i < len(values); i++ {
		screen.draw(i, values[i])
	}
}
