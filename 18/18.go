package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"slices"
)

type Cube struct {
	X, Y, Z int
}

func load[T any](filename string) []T {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var arr []T
	err = json.Unmarshal(data, &arr)
	if err != nil {
		panic(err)
	}

	return arr
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func Neighbors(cube Cube) []Cube {
	return []Cube{
		{cube.X - 1, cube.Y, cube.Z},
		{cube.X + 1, cube.Y, cube.Z},
		{cube.X, cube.Y - 1, cube.Z},
		{cube.X, cube.Y + 1, cube.Z},
		{cube.X, cube.Y, cube.Z - 1},
		{cube.X, cube.Y, cube.Z + 1},
	}
}

func solve(cubes []Cube) {
	positions := make(map[Cube]bool)
	for _, cube := range cubes {
		positions[cube] = true
	}
	totalSurface := 0
	for _, cube := range cubes {
		surface := 0
		for _, c := range Neighbors(cube) {
			if _, ok := positions[c]; !ok {
				surface += 1
			}
		}
		totalSurface += surface
	}
	fmt.Println("A:", totalSurface)

	var minCube, maxCube Cube
	for _, cube := range cubes {
		minCube.X = min(minCube.X, cube.X)
		minCube.Y = min(minCube.Y, cube.Y)
		minCube.Z = min(minCube.Z, cube.Z)
		maxCube.X = max(maxCube.X, cube.X)
		maxCube.Y = max(maxCube.Y, cube.Y)
		maxCube.Z = max(maxCube.Z, cube.Z)
	}
	minCube.X--
	minCube.Y--
	minCube.Z--
	maxCube.X++
	maxCube.Y++
	maxCube.Z++
	fmt.Println(minCube, maxCube)

	seen := make(map[Cube]bool)
	q := make([]Cube, 0)
	q = append(q, minCube)
	seen[minCube] = true
	totalSurface = 0
	for len(q) > 0 {
		cube := q[0]
		q = q[1:]
		neighbors := Neighbors(cube)
		// remove out of bounds
		neighbors = slices.DeleteFunc(neighbors, func(c Cube) bool {
			return c.X < minCube.X || c.X > maxCube.X ||
				c.Y < minCube.Y || c.Y > maxCube.Y ||
				c.Z < minCube.Z || c.Z > maxCube.Z
		})
		for _, c := range neighbors {
			_, seenBefore := seen[c]
			seen[c] = true
			_, isOccupied := positions[c]
			if isOccupied {
				totalSurface += 1
			} else if !seenBefore {
				q = append(q, c)
			}
		}
	}

	fmt.Println("B:", totalSurface)
}

func main() {
	entries, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	idx := slices.IndexFunc[[]fs.DirEntry](entries, func(e fs.DirEntry) bool { return e.Name() == "ex.json" })
	if idx == -1 {
		log.Fatal("Could not find example file ex.json")
	}
	idx = slices.IndexFunc[[]fs.DirEntry](entries, func(e fs.DirEntry) bool { return e.Name() == "in.json" })
	if idx == -1 {
		log.Fatal("Could not find input file in.json")
	}
	fmt.Println("------")
	fmt.Println("Example")
	fmt.Println("------")
	cubes := load[Cube]("ex.json")
	solve(cubes)
	fmt.Println("------")
	fmt.Println("Test")
	fmt.Println("------")
	cubes = load[Cube]("in.json")
	solve(cubes)
}
