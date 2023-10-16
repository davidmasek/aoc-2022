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

func sum1D(arr []bool) int {
	s := 0
	for _, v := range arr {
		if v {
			s++
		}
	}
	return s
}

func sum2D(arr [][]bool) int {
	s := 0
	for _, v := range arr {
		s += sum1D(v)
	}
	return s
}

func traverse(maze [][]int, currentHeight int, r, c int, dx, dy int, rowCount, colCount int) int {
	// literal edge cases
	if r == 0 || c == 0 || r == rowCount-1 || c == colCount-1 {
		return 0
	}
	sum := 0
	for !(r == 0 || c == 0 || r == rowCount-1 || c == colCount-1) {
		r += dy
		c += dx
		sum++
		if maze[r][c] >= currentHeight {
			break
		}
	}
	return sum
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

	maze := make([][]int, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		row := make([]int, 0)
		for _, c := range line {
			v := int(c - '0')
			row = append(row, v)
		}
		maze = append(maze, row)
	}

	h, w := len(maze), len(maze[0])
	visible := make([][]bool, h)
	for i := range visible {
		visible[i] = make([]bool, w)
	}
	if !test {
		fmt.Println(h, w)
		fmt.Println(maze)
	}
	for r := range visible {
		// from left
		currentMaxHeight := -1
		for c := range visible[r] {
			if maze[r][c] > currentMaxHeight {
				visible[r][c] = true
			}
			if maze[r][c] > currentMaxHeight {
				currentMaxHeight = maze[r][c]
			}
		}
		// from right
		currentMaxHeight = -1
		for c := len(visible[r]) - 1; c >= 0; c-- {
			if maze[r][c] > currentMaxHeight {
				visible[r][c] = true
			}
			if maze[r][c] > currentMaxHeight {
				currentMaxHeight = maze[r][c]
			}
		}
	}
	for c := range visible[0] {
		// from top
		currentMaxHeight := -1
		for r := range visible {
			if maze[r][c] > currentMaxHeight {
				visible[r][c] = true
			}
			if maze[r][c] > currentMaxHeight {
				currentMaxHeight = maze[r][c]
			}
		}
		// from bottom
		currentMaxHeight = -1
		for r := len(visible) - 1; r >= 0; r-- {
			if maze[r][c] > currentMaxHeight {
				visible[r][c] = true
			}
			if maze[r][c] > currentMaxHeight {
				currentMaxHeight = maze[r][c]
			}
		}
	}
	if !test {
		fmt.Println(visible)
	}

	fmt.Println("A:", sum2D(visible))

	// part B
	maxScore := 0
	i, j := 0, 0
	for r := range visible {
		for c := range visible[r] {
			score := 1
			score = score * traverse(maze, maze[r][c], r, c, 1, 0, h, w)
			score = score * traverse(maze, maze[r][c], r, c, -1, 0, h, w)
			score = score * traverse(maze, maze[r][c], r, c, 0, 1, h, w)
			score = score * traverse(maze, maze[r][c], r, c, 0, -1, h, w)
			if score > maxScore {
				maxScore = score
				i, j = r, c
			}
		}
	}
	fmt.Println("B:", maxScore, i, j)
}
