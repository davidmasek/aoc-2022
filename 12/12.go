package main

import (
	"bufio"
	"fmt"
	"os"
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

func InBounds(r, c, height, width int) bool {
	return r >= 0 && r < height && c >= 0 && c < width
}

func Neighbors(n Node, height, width int, visited [][]bool, maze [][]int) []Node {
	available := make([]Node, 0)
	for _, nb := range []Node{{n.r - 1, n.c}, {n.r + 1, n.c}, {n.r, n.c - 1}, {n.r, n.c + 1}} {
		if InBounds(nb.r, nb.c, height, width) && !visited[nb.r][nb.c] {
			if B_SIDE && maze[n.r][n.c] <= maze[nb.r][nb.c]+1 {
				available = append(available, nb)
			} else if !B_SIDE && maze[n.r][n.c]+1 >= maze[nb.r][nb.c] {
				available = append(available, nb)
			}
		}
	}
	return available
}

type Node struct {
	r, c int
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
	var maze [][]int = make([][]int, 0)
	i := 0
	var start, goal Node
	for scanner.Scan() {
		line := scanner.Text()
		var row []int = make([]int, 0)
		for j, c := range line {
			if c == 'S' {
				row = append(row, 0)
				start.r = i
				start.c = j
			} else if c == 'E' {
				row = append(row, int('z'-'a'))
				goal.r = i
				goal.c = j
			} else {
				row = append(row, int(c-'a'))
			}
		}
		maze = append(maze, row)
		i++
	}

	height, width := len(maze), len(maze[0])

	var queue []Node = make([]Node, 0)
	var visited [][]bool = make([][]bool, height)
	for i := 0; i < height; i++ {
		visited[i] = make([]bool, width)
	}

	distances := make([][]int, height)
	for i := 0; i < height; i++ {
		distances[i] = make([]int, width)
		for j := 0; j < width; j++ {
			distances[i][j] = -1
		}
	}
	// start
	if B_SIDE {
		queue = append(queue, Node{goal.r, goal.c})
	} else {
		queue = append(queue, Node{start.r, start.c})
	}
	n := queue[0]
	visited[n.r][n.c] = true
	var prev map[Node]Node = make(map[Node]Node)
	distances[n.r][n.c] = 0

	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		if !B_SIDE && n.r == goal.r && n.c == goal.c {
			fmt.Println("Found path")
			break
		}
		for _, nb := range Neighbors(n, height, width, visited, maze) {
			queue = append(queue, nb)
			visited[nb.r][nb.c] = true
			distances[nb.r][nb.c] = distances[n.r][n.c] + 1
			prev[nb] = n
		}
	}
	if !TEST {
		// pretty print distances from start
		for i := 0; i < height; i++ {
			for j := 0; j < width; j++ {
				fmt.Printf("%2d ", distances[i][j])
			}
			fmt.Println()
		}
		fmt.Println("-----")
	}
	// print only those whose height is zero on single lines
	min := distances[0][0]
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if maze[i][j] == 0 {
				if !TEST {
					fmt.Printf("(%d,%d): (%d)\n", i, j, distances[i][j])
				}
				if distances[i][j] < min && distances[i][j] != -1 {
					min = distances[i][j]
				}
			}
		}
	}
	if !TEST {
		fmt.Println("-----")
	}
	fmt.Println("from start", distances[start.r][start.c])
	fmt.Println("nearest", min)
}
