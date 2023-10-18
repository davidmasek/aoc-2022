package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Point struct {
	x, y  int
	label string
}

type Rope = []Point

type Maze struct {
	rope []Point
}

func NewMaze(ropeLength int) *Maze {
	if ropeLength < 2 {
		panic("Rope too short")
	}
	if ropeLength > 10 {
		panic("Rope too long")
	}
	m := new(Maze)
	m.rope = make([]Point, ropeLength)
	for i := 0; i < ropeLength; i++ {
		m.rope[i].label = fmt.Sprint(i)
	}
	m.rope[0].label = "H"
	m.rope[ropeLength-1].label = "T"
	return m
}

func (m *Maze) head() *Point {
	return &m.rope[0]
}

func (m *Maze) tail() *Point {
	return &m.rope[len(m.rope)-1]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}

func signClip(x int, absSize int) int {
	if x > absSize {
		return absSize
	}
	if x < -absSize {
		return -absSize
	}
	return x
}

func (m *Maze) render() {
	start_y := min(m.head().y, m.tail().y)
	start_x := min(m.head().x, m.tail().x)
	fmt.Println("-----------")
	for y := start_y; y < start_y+4; y++ {
		fmt.Print("|")
		for x := start_x; x < start_x+4; x++ {
			occupied := false
			for _, p := range m.rope {
				if x == p.x && y == p.y {
					fmt.Print(p.label)
					occupied = true
					break
				}
			}
			if !occupied {
				fmt.Print(".")
			}
		}
		fmt.Printf("  %3d|\n", y)
	}
	fmt.Printf("|%-3d %-5d|\n", start_x, start_x+3)
	fmt.Println("-----------")
}

func snap(head, tail *Point) {
	if abs(head.x-tail.x) > 1 || abs(head.y-tail.y) > 1 {
		tail.x += signClip(head.x-tail.x, 1)
		tail.y += signClip(head.y-tail.y, 1)
	}
}

func (m *Maze) move(dx, dy int) {
	m.head().x += dx
	m.head().y += dy
	for i := 0; i < len(m.rope)-1; i++ {
		snap(&m.rope[i], &m.rope[i+1])
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

	var mazeLength int
	const VAR_B = true
	if VAR_B {
		mazeLength = 10
	} else {
		mazeLength = 2
	}

	var m Maze = *NewMaze(mazeLength)
	if !test {
		m.render()
	}
	visited := mapset.NewSet[Point]()
	visited.Add(*m.tail())

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		var dir string
		var dist int
		fmt.Sscanf(line, "%s %d", &dir, &dist)
		for i := 0; i < dist; i++ {
			switch dir {
			case "U":
				m.move(0, -1)
			case "D":
				m.move(0, 1)
			case "L":
				m.move(-1, 0)
			case "R":
				m.move(1, 0)
			default:
				fmt.Println("Unknown direction")
				continue
			}
			if !test {
				fmt.Println(dir)
				m.render()
			}
			visited.Add(*m.tail())
		}
	}
	fmt.Println(visited.Cardinality())
}
