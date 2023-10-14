package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Range struct {
	lower int
	upper int
}

func (r *Range) contains(o *Range) bool {
	return r.lower <= o.lower && r.upper >= o.upper
}

func (r *Range) overlaps(o *Range) bool {
	return r.upper >= o.lower && r.lower <= o.upper
}

func rangeFromString(s string) Range {
	res := strings.Split(s, "-")
	lower, err := strconv.Atoi(res[0])
	check(err)
	upper, err := strconv.Atoi(res[1])
	check(err)
	return Range{lower: lower, upper: upper}
}

func main() {
	fmt.Println(len(os.Args), os.Args)
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

	var containsCount int = 0
	var overlapsCount int = 0
	for {
		var line string
		n, _ := fmt.Fscanln(f, &line)
		if n != 1 {
			break
		}
		res := strings.Split(line, ",")
		left := res[0]
		right := res[1]
		var first Range = rangeFromString(left)
		var second Range = rangeFromString(right)
		if first.contains(&second) || second.contains(&first) {
			containsCount++
		}
		if first.overlaps(&second) {
			overlapsCount++
		}
	}
	fmt.Println("A:", containsCount)
	fmt.Println("B:", overlapsCount)
}
