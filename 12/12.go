package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const LOGGING = false
const TEST = true
const B_SIDE = false

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
		if InBounds(nb.r, nb.c, height, width) && !visited[nb.r][nb.c] && maze[n.r][n.c]+1 >= maze[nb.r][nb.c] {
			available = append(available, nb)
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
	// start
	queue = append(queue, Node{start.r, start.c})
	visited[start.r][start.c] = true
	var prev map[Node]Node = make(map[Node]Node)

	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		if n.r == goal.r && n.c == goal.c {
			fmt.Println("Found path")
			break
		}
		for _, nb := range Neighbors(n, height, width, visited, maze) {
			queue = append(queue, nb)
			visited[nb.r][nb.c] = true
			prev[nb] = n
		}
	}
	path := make([]Node, 0)
	c := goal
	path = append(path, c)
	for {
		c = prev[c]
		path = append(path, c)
		if c == start {
			break
		}
	}
	for i := len(path) - 1; i >= 0; i-- {
		debug(path[i])
	}
	fmt.Println(len(path) - 1)
}
