package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const LOGGING = false
const TEST = true
const B_SIDE = true

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func debug(args ...interface{}) {
	if LOGGING {
		fmt.Println(args...)
	}
}

type Point struct {
	x, y int
}

type Blocker interface {
	Contains(p Point) bool
}

// Vertical/horizontal line
type Line struct {
	p1, p2 Point
}

type FloorLine struct {
	y int
}

func (line Line) Contains(p Point) bool {
	return p.x >= line.p1.x && p.x <= line.p2.x && p.y >= line.p1.y && p.y <= line.p2.y ||
		p.x >= line.p2.x && p.x <= line.p1.x && p.y >= line.p2.y && p.y <= line.p1.y
}

func (line FloorLine) Contains(p Point) bool {
	return p.y == line.y
}

func Blocked(point Point, lines []Blocker) bool {
	for _, line := range lines {
		if line.Contains(point) {
			return true
		}
	}
	return false
}

type SandPile struct {
	pile map[Point]bool
}

func (pile SandPile) Contains(p Point) bool {
	_, ok := pile.pile[p]
	return ok
}

func main() {
	fmt.Println(len(os.Args), os.Args)
	fmt.Println("-----")
	if len(os.Args) != 2 {
		fmt.Println("Usage: XY.go XY.in")
		os.Exit(1)
	}
	var filename string
	if TEST {
		filename = strings.Replace(os.Args[1], ".ex", ".in", 1)
	} else {
		filename = os.Args[1]
	}
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// reader := bufio.NewReader(f)

	linesCount := 0
	lines := make([]Blocker, 0)
	maxY := 0
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Println(line)
		out := strings.Split(line, "->")
		for i := 0; i < len(out)-1; i++ {
			out[i] = strings.TrimSpace(out[i])
			// fmt.Println(out[i], out[i+1])
			out[i+1] = strings.TrimSpace(out[i+1])
			xyFrom := strings.Split(out[i], ",")
			xyTo := strings.Split(out[i+1], ",")
			if len(xyFrom) != 2 || len(xyTo) != 2 {
				fmt.Println("Bad line:", line)
				os.Exit(1)
			}
			// fmt.Println(xyFrom, xyTo)
			xFrom, err := strconv.Atoi(xyFrom[0])
			check(err)
			yFrom, err := strconv.Atoi(xyFrom[1])
			if yFrom > maxY {
				maxY = yFrom
			}
			check(err)
			xTo, err := strconv.Atoi(xyTo[0])
			check(err)
			yTo, err := strconv.Atoi(xyTo[1])
			if yTo > maxY {
				maxY = yTo
			}
			check(err)
			lines = append(lines, Line{Point{xFrom, yFrom}, Point{xTo, yTo}})
		}
		linesCount += len(out)
	}

	if B_SIDE {
		lines = append(lines, FloorLine{maxY + 2})
	}
	source := Point{500, 0}

	sand := Point{source.x, source.y}
	sandCount := 0
	var pile SandPile = SandPile{make(map[Point]bool)}
	lines = append(lines, pile)

	for {
		if sand.y > 1000 {
			fmt.Println("abyss ~~")
			break
		}
		if !Blocked(Point{sand.x, sand.y + 1}, lines) {
			sand.y++
			continue
		}
		if !Blocked(Point{sand.x - 1, sand.y + 1}, lines) {
			sand.x--
			sand.y++
			continue
		}
		if !Blocked(Point{sand.x + 1, sand.y + 1}, lines) {
			sand.x++
			sand.y++
			continue
		}
		sandCount++
		pile.pile[sand] = true
		// lines = append(lines, Line{Point{sand.x, sand.y}, Point{sand.x, sand.y}})
		if sand.x == source.x && sand.y == source.y {
			fmt.Println("src blocked")
			break
		}
		sand = Point{source.x, source.y}
	}
	fmt.Println("maxY", maxY)
	fmt.Println("Sand count:", sandCount)
}
