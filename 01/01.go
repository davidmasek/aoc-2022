package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println(len(os.Args), os.Args)
	if len(os.Args) != 2 {
		fmt.Println("Usage: XY.go XY.in")
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	check(err)
	defer f.Close()

	// var lines []string
	var currentSum int = 0
	var allValues []int
	for {
		var line string
		n, err := fmt.Fscanln(f, &line)
		if n > 1 {
			panic("too many values")
		}
		if n == 1 {
			v, err := strconv.Atoi(line)
			check(err)
			currentSum += v
			// lines = append(lines, line)
		}
		if n == 0 {
			allValues = append(allValues, currentSum)
			currentSum = 0
		}
		if err == io.EOF {
			break
		}
	}
	slices.Sort(allValues)
	var sum = 0
	for _, value := range allValues[len(allValues)-3:] {
		sum += value
	}
	fmt.Println(sum)
}
