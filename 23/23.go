package main

import (
	"fmt"
	"strings"
)

const EX = `
....#..
..###.#
#...#.#
.#...##
#.###..
##.#.##
.#..#..
`

const EX_SMALL = `
.....
..##.
..#..
.....
..##.
.....
`

type Point struct {
	X, Y int
}

type MultiPoint struct {
	X, Y, Count int
}

func load(lines string) ([][]bool, map[Point]bool) {
	rows := strings.Split(lines, "\n")
	res := make([][]bool, len(rows))
	points := make(map[Point]bool, 0)
	for y, row := range rows {
		res[y] = make([]bool, len(row))
		for x, ch := range row {
			res[y][x] = ch == '#'
			if ch == '#' {
				points[Point{X: x, Y: y}] = true
			}
		}
	}
	return res, points
}

func PrintGrid(points map[Point]bool, proposed map[Point]MultiPoint) {
	return
	minX, maxX, minY, maxY := getRange(points)
	if len(proposed) > 0 {
		minXp, maxXp, minYp, maxYp := getRange(proposed)
		minX = min(minX, minXp)
		maxX = max(maxX, maxXp)
		minY = min(minY, minYp)
		maxY = max(maxY, maxYp)
	}

	fmt.Println("[Y,X]:", minY, minX)
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if p, o := proposed[Point{X: x, Y: y}]; o {
				fmt.Print(p.Count)
			} else if points[Point{X: x, Y: y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println("[Y,X]:", maxY, maxX)

	if len(proposed) > 0 {
		fmt.Println("|-|-|-|-|-|-|----------|")
	} else {
		fmt.Println("|------------|-|-|-|-|-|")
	}
}

func (p Point) Add(o Point) Point {
	return Point{X: p.X + o.X, Y: p.Y + o.Y}
}

func (p MultiPoint) Add(o Point) MultiPoint {
	return MultiPoint{X: p.X + o.X, Y: p.Y + o.Y, Count: p.Count}
}

func getRange[T any](points map[Point]T) (minX, maxX, minY, maxY int) {
	first := true
	for p := range points {
		if first {
			minX, maxX, minY, maxY = p.X, p.X, p.Y, p.Y
			first = false
			continue
		}
		minX = min(minX, p.X)
		maxX = max(maxX, p.X)
		minY = min(minY, p.Y)
		maxY = max(maxY, p.Y)
	}
	return
}

func ProposeMove(elves map[Point]bool, elf Point, directions [][]Point) Point {
	nextPossible := make([]Point, 0)
	for _, dir := range directions {
		canMove := true
		for _, move := range dir {
			_, occupied := elves[Point{elf.X + move.X, elf.Y + move.Y}]
			if occupied {
				canMove = false
				break
			}
		}
		if canMove {
			nextPos := elf.Add(dir[1])
			nextPossible = append(nextPossible, nextPos)
		}
	}
	// if cannot move anywhere -> stay
	// if all around empty -> also stay
	if len(nextPossible) == 0 || len(nextPossible) == 4 {
		return elf
	}
	// else choose first available (order is important)
	return nextPossible[0]
}

func solve(elves map[Point]bool) {
	PrintGrid(elves, make(map[Point]MultiPoint))

	// k := []string{"N", "S", "W", "E"}
	// fmt.Println(k)
	// k = append(k[1:], k[0])
	// fmt.Println(k)
	// k = append(k[1:], k[0])
	// fmt.Println(k)
	// k = append(k[1:], k[0])
	// fmt.Println(k)
	// return

	// N,S,W,E - check diagonals but move only N,S,W,E
	// first propose moves and if possible (no collision) then move
	// then the checked order rotates (N,S,W,E -> S,W,E,N)
	directions := [][]Point{
		{{-1, -1}, {0, -1}, {1, -1}}, // N
		{{-1, 1}, {0, 1}, {1, 1}},    // S
		{{-1, -1}, {-1, 0}, {-1, 1}}, // W
		{{1, -1}, {1, 0}, {1, 1}},    // E
	}

	for i := 0; i < 1000; i++ {

		// make propositions
		proposed := make(map[Point]MultiPoint, 0)
		for elf := range elves {
			nextPos := ProposeMove(elves, elf, directions)
			_, alreadyProposed := proposed[nextPos]
			if alreadyProposed {
				pt := proposed[nextPos]
				pt.Count++
				proposed[nextPos] = pt
			} else {
				proposed[nextPos] = MultiPoint{nextPos.X, nextPos.Y, 1}
			}
		}
		PrintGrid(elves, proposed)
		// move
		newPositions := make(map[Point]bool, 0)
		for elf := range elves {
			nextPos := ProposeMove(elves, elf, directions)
			proposed, found := proposed[nextPos]
			if !found {
				panic("No proposed move")
			}
			// keep at old position
			if proposed.Count > 1 {
				newPositions[elf] = true
			} else {
				// move to new position
				newPositions[Point{proposed.X, proposed.Y}] = true
			}
		}
		fmt.Println("End of round", i+1)
		PrintGrid(newPositions, make(map[Point]MultiPoint))
		same := true
		for elf := range elves {
			_, found := newPositions[elf]
			if !found {
				same = false
				break
			}
		}
		elves = newPositions
		// rotate directions
		directions = append(directions[1:], directions[0])

		// stop on movement
		if same {
			fmt.Println("Stopping after round", i+1)
			fmt.Println("B:", i+1)
			break
		}
	}

	minX, maxX, minY, maxY := getRange(elves)
	emptyCount := 0
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if _, occ := elves[Point{X: x, Y: y}]; !occ {
				emptyCount++
			}
		}
	}
	fmt.Println("A:", emptyCount)
}

func main() {
	// _, elves := load(EX)
	// solve(elves)
	_, elves := load(INPUT)
	solve(elves)
}
